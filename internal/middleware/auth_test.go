package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	middlewareFunc := AuthMiddleware(nil)
	middlewareFunc(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.True(t, c.IsAborted())
}

func TestAuthMiddleware_InvalidFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "InvalidFormat")

	middlewareFunc := AuthMiddleware(nil)
	middlewareFunc(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.True(t, c.IsAborted())
}

func TestAuthMiddleware_NotBearer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Basic token123")

	middlewareFunc := AuthMiddleware(nil)
	middlewareFunc(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.True(t, c.IsAborted())
}

func TestGetUserID_NotSet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	userID := GetUserID(c)

	assert.Equal(t, uint(0), userID)
}

func TestGetUserID_Set(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", uint(123))

	userID := GetUserID(c)

	assert.Equal(t, uint(123), userID)
}

func TestGetUserID_WrongType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", "not_a_uint")

	userID := GetUserID(c)

	assert.Equal(t, uint(0), userID)
}

func TestGetUserEmail_NotSet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	email := GetUserEmail(c)

	assert.Equal(t, "", email)
}

func TestGetUserEmail_Set(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("email", "test@example.com")

	email := GetUserEmail(c)

	assert.Equal(t, "test@example.com", email)
}

func TestGetUserEmail_WrongType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("email", 123)

	email := GetUserEmail(c)

	assert.Equal(t, "", email)
}

func TestGetUserRole_NotSet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	role := GetUserRole(c)

	assert.Equal(t, "", role)
}

func TestGetUserRole_Set(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "admin")

	role := GetUserRole(c)

	assert.Equal(t, "admin", role)
}

func TestGetUserRole_WrongType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", 456)

	role := GetUserRole(c)

	assert.Equal(t, "", role)
}
