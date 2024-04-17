package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sahildhargave/memories/account/model"
	"github.com/sahildhargave/memories/account/model/apperrors"
	"github.com/sahildhargave/memories/account/model/mocks"
	"github.com/stretchr/testify/assert"
)

func TestAuthUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a new instance of the mock token service
	mockTokenService := new(mocks.MockTokenService)

	// Generate a random UUID for testing
	uid := uuid.New()

	// Create a mock user for testing
	u := &model.User{
		UID:   uid,
		Email: "sam@gmail.com",
	}

	// Define valid and invalid token strings
	validTokenHeader := "validTokenString"
	invalidTokenHeader := "invalidTokenString"
	invalidTokenErr := apperrors.NewAuthorization("Unable to verify user from idToken")

	// Configure mockTokenService expectations
	mockTokenService.On("ValidateIDToken", validTokenHeader).Return(u, nil)
	mockTokenService.On("ValidateIDToken", invalidTokenHeader).Return(nil, invalidTokenErr)

	t.Run("Adds a user to context with valid token", func(t *testing.T) {
		rr := httptest.NewRecorder()
		r := gin.New()

		r.GET("/me", AuthUser(mockTokenService), func(c *gin.Context) {
			contextUser, exists := c.Get("user")
			assert.True(t, exists)
			assert.Equal(t, u, contextUser)
		})

		request, _ := http.NewRequest(http.MethodGet, "/me", nil)
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", validTokenHeader))

		r.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Handles invalid token", func(t *testing.T) {
		rr := httptest.NewRecorder()
		r := gin.New()

		r.GET("/me", AuthUser(mockTokenService))

		request, _ := http.NewRequest(http.MethodGet, "/me", nil)
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", invalidTokenHeader))

		r.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Handles missing authorization header", func(t *testing.T) {
		rr := httptest.NewRecorder()
		r := gin.New()

		r.GET("/me", AuthUser(mockTokenService))

		request, _ := http.NewRequest(http.MethodGet, "/me", nil)

		r.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})
}
