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
	// Setup
	gin.SetMode(gin.TestMode)

	t.Run("Email and Password Required", func(t *testing.T) {
		// We just want this to show that it's not called in this case
		mockUserService := new(mocks.MockUserService)
		mockUserService.On("Signup", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*model.User")).Return(nil)

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
			"email": "",
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

	t.Run("Invalid email", func(t *testing.T) {
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
			"email":    "bob@bob",
			"password": "supersecret1234",
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

	t.Run("Password too short", func(t *testing.T) {
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
			"email":    "bob@bob.com",
			"password": "supe",
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
			"email":    "bob@bob.com",
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
			Email:    "bob@bob.com",
			Password: "avalidpassword",
		}

		mockUserService := new(mocks.MockUserService)
		mockUserService.On("Signup", mock.Anything, u).
			Return(apperrors.NewConflict("User Already Exists", u.Email))

		rr := httptest.NewRecorder()
		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		// Create a request body with valid email and password
		reqBody, err := json.Marshal(gin.H{
			"email":    u.Email,
			"password": u.Password,
		})
		assert.NoError(t, err)

		// Send a POST request to "/signup" endpoint with the request body
		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		// Assert that the response status code matches the expected 409 Conflict
		assert.Equal(t, http.StatusConflict, rr.Code)

		// Assert that the UserService.Signup method was called with the expected parameters
		mockUserService.AssertExpectations(t)
	})

t.Run("Successful Token Creation", func(t *testing.T) {
    // Create a new Gin router
    router := gin.Default()

    // Create a mock user and token response
    u := &model.User{
        Email:    "bob@bob.com",
        Password: "avalidpassword",
    }
    mockTokenResp := &model.TokenPair{
        IDToken:      "idToken",
        RefreshToken: "refreshToken",
    }

    // Create mock services
    mockUserService := &mocks.MockUserService{}
    mockTokenService := &mocks.MockTokenService{}

    // Configure mock service responses
    mockUserService.On("Signup", mock.Anything, u).Return(nil)
    mockTokenService.On("NewPairFromUser", mock.Anything, u, "").Return(mockTokenResp, nil)

    // Create a handler instance with mock services injected
    NewHandler(&Config{
        R:            router,
        UserService:  mockUserService,
        TokenService: mockTokenService,
    })

    // Prepare request body
    reqBody, err := json.Marshal(gin.H{
        "email":    u.Email,
        "password": u.Password,
    })
    assert.NoError(t, err)

    // Create HTTP request
    request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
    assert.NoError(t, err)
    request.Header.Set("Content-Type", "application/json")

    // Create a new response recorder to capture the HTTP response
    rr := httptest.NewRecorder()

    // Perform the HTTP request
    router.ServeHTTP(rr, request)

    // Assert the expected response status code
    assert.Equal(t, http.StatusCreated, rr.Code)

    // Assert the expected response body
    var respBody map[string]interface{}
    err = json.Unmarshal(rr.Body.Bytes(), &respBody)
    assert.NoError(t, err)

    expectedRespBody := map[string]interface{}{
        "tokens": mockTokenResp,
    }
    assert.Equal(t, expectedRespBody, respBody)

    // Assert that the expected methods were called on the mock services
    mockUserService.AssertExpectations(t)
    mockTokenService.AssertExpectations(t)
})

	t.Run("Failed Token Creation", func(t *testing.T) {
    // Create a mock user and error response
    u := &model.User{
        Email:    "bob@bob.com",
        Password: "avalidpassword",
    }
    mockErrorResponse := apperrors.NewInternal()

    // Create mock services
    mockUserService := &mocks.MockUserService{}
    mockTokenService := &mocks.MockTokenService{}

    // Configure mock service responses
    mockUserService.On("Signup", mock.Anything, u).Return(nil)
    mockTokenService.On("NewPairFromUser", mock.Anything, u, "").Return(nil, mockErrorResponse)

    // Create a new Gin router
    router := gin.Default()

    // Create a handler instance with mock services injected
    NewHandler(&Config{
        R:            router,
        UserService:  mockUserService,
        TokenService: mockTokenService,
    })

    // Prepare request body
    reqBody, err := json.Marshal(gin.H{
        "email":    u.Email,
        "password": u.Password,
    })
    assert.NoError(t, err)

    // Create HTTP request
    request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
    assert.NoError(t, err)
    request.Header.Set("Content-Type", "application/json")

    // Create a new response recorder to capture the HTTP response
    rr := httptest.NewRecorder()

    // Perform the HTTP request
    router.ServeHTTP(rr, request)

    // Assert the expected response status code (internal server error)
    assert.Equal(t, http.StatusInternalServerError, rr.Code)

    // Assert the expected response body
    var respBody map[string]interface{}
    err = json.Unmarshal(rr.Body.Bytes(), &respBody)
    assert.NoError(t, err)

    expectedRespBody := map[string]interface{}{
        "error": mockErrorResponse,
    }
    assert.Equal(t, expectedRespBody, respBody)

    // Assert that the expected methods were called on the mock services
    mockUserService.AssertExpectations(t)
    mockTokenService.AssertExpectations(t)
})
}

