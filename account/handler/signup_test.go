package handler

import (
	"bytes"
	"encoding/json"

	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sahildhargave/memories/account/model"
	"github.com/sahildhargave/memories/account/model/apperrors"
	"github.com/sahildhargave/memories/account/model/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignup(t *testing.T) {
	//	TODO// set
	gin.SetMode(gin.TestMode)

	t.Run("Email and Password Required", func(t *testing.T) {

		mockUserService := new(mocks.MockUserService)
		mockUserService.On("signup", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*model.User")).Return(nil)

		// aresponse recorder for getting written response

		rr := httptest.NewRecorder()

		// not middlerware as we dont't yet have authorized user
		router := gin.Default()
		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		// create a request body with empty email and password

		reqBody, err := json.Marshal(gin.H{
			"email": "",
		})

		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		// TODO bytes.NewBuffer to create a reader

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		assert.Equal(t, 400, rr.Code)
		mockUserService.AssertNotCalled(t, "Signup")
	})

	t.Run("Invalid email", func(t *testing.T) {

		mockUserService := new(mocks.MockUserService)
		mockUserService.On("signup", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*model.User")).Return(nil)

		rr := httptest.NewRecorder()

		// need no middleware as we dont yet have authorized user
		router := gin.Default()

		NewHandler(&Config{R: router,
			UserService: mockUserService,
		})

		reqBody, err := json.Marshal(gin.H{
			"email":    "sam@gmail",
			"password": "123456",
		})

		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, 400, rr.Code)
		mockUserService.AssertNotCalled(t, "Signup")
	})

	t.Run("Password Too Short", func(t *testing.T) {

		mockUserService := new(mocks.MockUserService)
		mockUserService.On("Signup", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*model.User")).Return(nil)

		// a response recorder
		rr := httptest.NewRecorder()

		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		reqBody, err := json.Marshal(gin.H{
			"email":    "sam@gmail.com",
			"password": "123",
		})

		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, 400, rr.Code)
		mockUserService.AssertNotCalled(t, "Signup")
	})

	t.Run("Password too long", func(t *testing.T) {
		// We just want this to show that it's not called in this case
		mockUserService := new(mocks.MockUserService)
		mockUserService.On("Signup", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*model.User")).Return(nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// don't need a middleware as we don't yet have authorized user
		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		// create a request body with empty email and password
		reqBody, err := json.Marshal(gin.H{
			"email":    "sam@gamil.com",
			"password": "super12324jhklafsdjhflkjweyruasdljkfhasdldfjkhasdkljhrleqwwjkrhlqwejrhasdflkjhasdf",
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, 400, rr.Code)
		mockUserService.AssertNotCalled(t, "Signup")
	})
	t.Run("Error returned from UserService", func(t *testing.T) {
		u := &model.User{
			Email:    "sam@gmail.com",
			Password: "123123",
		}

		mockUserService := new(mocks.MockUserService)
		mockUserService.On("Signup", mock.AnythingOfType("*gin.Context"), u).Return(apperrors.NewConflict("User Already Exists", u.Email))

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// don't need a middleware as we don't yet have authorized user
		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		// create a request body with empty email and password
		reqBody, err := json.Marshal(gin.H{
			"email":    u.Email,
			"password": u.Password,
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, 409, rr.Code)
		mockUserService.AssertExpectations(t)
	})

	t.Run("Successful Token Creation", func(t *testing.T) {
		u := &model.User{
			Email:    "bob@bob.com",
			Password: "avalidpassword",
		}
		mockTokenResp := &model.TokenPair{
			IDToken:      "idToken",
			RefreshToken: "refreshToken",
		}

		// Mock UserService and TokenService
		mockUserService := new(mocks.MockUserService)
		mockTokenService := new(mocks.MockTokenService)

		// Configure expected calls and return values for UserService and TokenService
		mockUserService.On("Signup", mock.AnythingOfType("*gin.Context"), u).Return(nil)
		mockTokenService.On("NewPairFromUser", mock.AnythingOfType("*gin.Context"), u, "").Return(mockTokenResp, nil)

		// Create a response recorder to capture HTTP response
		rr := httptest.NewRecorder()

		// Create a Gin router
		router := gin.Default()

		// Inject mocked services into handler
		NewHandler(&Config{
			R:            router,
			UserService:  mockUserService,
			TokenService: mockTokenService,
		})

		// Create request body as JSON
		reqBody, err := json.Marshal(gin.H{
			"email":    u.Email,
			"password": u.Password,
		})
		assert.NoError(t, err)

		// Create HTTP request with the request body
		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")

		// Perform HTTP request
		router.ServeHTTP(rr, request)

		// Validate HTTP response
		assert.Equal(t, http.StatusCreated, rr.Code, "Expected HTTP status 201")

		// Validate response body
		//respBody, err := json.Marshal(gin.H{"tokens": mockTokenResp})
		//assert.NoError(t, err)
		//assert.JSONEq(t, string(respBody), rr.Body.String(), "Response body does not match expected")

		// Assert expected method calls on mock services
		mockUserService.AssertExpectations(t)
		mockTokenService.AssertExpectations(t)
	})

	t.Run("Failed Token Creation", func(t *testing.T) {
		u := &model.User{
			Email:    "bob@bob.com",
			Password: "avalidpassword",
		}

		mockErrorResponse := apperrors.NewInternal()

		mockUserService := new(mocks.MockUserService)
		mockTokenService := new(mocks.MockTokenService)

		mockUserService.
			On("Signup", mock.AnythingOfType("*gin.Context"), u).
			Return(nil)
		mockTokenService.
			On("NewPairFromUser", mock.AnythingOfType("*gin.Context"), u, "").
			Return(nil, mockErrorResponse)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// don't need a middleware as we don't yet have authorized user
		router := gin.Default()

		NewHandler(&Config{
			R:            router,
			UserService:  mockUserService,
			TokenService: mockTokenService,
		})

		// create a request body with empty email and password
		reqBody, err := json.Marshal(gin.H{
			"email":    u.Email,
			"password": u.Password,
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockErrorResponse,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockErrorResponse.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockUserService.AssertExpectations(t)
		mockTokenService.AssertExpectations(t)
	})

}
