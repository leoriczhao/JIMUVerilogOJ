package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestOptionalAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Helper to generate valid token
	generateToken := func() string {
		claims := jwt.MapClaims{
			"user_id":  float64(1),
			"username": "testuser",
			"role":     "student",
			"exp":      time.Now().Add(time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString(jwtSecret)
		return tokenString
	}

	t.Run("No Token - Should Continue as Anonymous", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/test", nil)

		handler := OptionalAuth()
		handler(c)

		// Should not abort
		assert.False(t, c.IsAborted())
		// Should not set user context
		_, exists := c.Get("user_id")
		assert.False(t, exists)
	})

	t.Run("Valid Token - Should Set User Context", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/test", nil)
		c.Request.Header.Set("Authorization", "Bearer "+generateToken())

		handler := OptionalAuth()
		handler(c)

		// Should not abort
		assert.False(t, c.IsAborted())
		// Should set user context
		userID, exists := c.Get("user_id")
		assert.True(t, exists)
		assert.Equal(t, uint(1), userID)

		username, exists := c.Get("username")
		assert.True(t, exists)
		assert.Equal(t, "testuser", username)

		role, exists := c.Get("role")
		assert.True(t, exists)
		assert.Equal(t, "student", role)
	})

	t.Run("Invalid Token - Should Continue as Anonymous", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/test", nil)
		c.Request.Header.Set("Authorization", "Bearer invalid-token")

		handler := OptionalAuth()
		handler(c)

		// Should not abort
		assert.False(t, c.IsAborted())
		// Should not set user context
		_, exists := c.Get("user_id")
		assert.False(t, exists)
	})

	t.Run("Malformed Header - Should Continue as Anonymous", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/test", nil)
		c.Request.Header.Set("Authorization", "NotBearer token")

		handler := OptionalAuth()
		handler(c)

		// Should not abort
		assert.False(t, c.IsAborted())
		// Should not set user context
		_, exists := c.Get("user_id")
		assert.False(t, exists)
	})
}

func TestAuthRequired(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Helper to generate valid token
	generateToken := func() string {
		claims := jwt.MapClaims{
			"user_id":  float64(1),
			"username": "testuser",
			"role":     "student",
			"exp":      time.Now().Add(time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString(jwtSecret)
		return tokenString
	}

	t.Run("No Token - Should Abort", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/test", nil)

		handler := AuthRequired()
		handler(c)

		// Should abort
		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Valid Token - Should Set User Context", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/test", nil)
		c.Request.Header.Set("Authorization", "Bearer "+generateToken())

		handler := AuthRequired()
		handler(c)

		// Should not abort
		assert.False(t, c.IsAborted())
		// Should set user context
		userID, exists := c.Get("user_id")
		assert.True(t, exists)
		assert.Equal(t, uint(1), userID)
	})

	t.Run("Invalid Token - Should Abort", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/test", nil)
		c.Request.Header.Set("Authorization", "Bearer invalid-token")

		handler := AuthRequired()
		handler(c)

		// Should abort
		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
