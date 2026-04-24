package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCORS_AllowedOrigin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Origin", "http://localhost:3000")

	corsFunc := CORS([]string{"http://localhost:3000"})
	corsFunc(c)

	assert.Equal(t, "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
}

func TestCORS_DisallowedOrigin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Origin", "http://evil.com")

	corsFunc := CORS([]string{"http://localhost:3000"})
	corsFunc(c)

	assert.Equal(t, "", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
}

func TestCORS_Wildcard(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Origin", "http://any.domain.com")

	corsFunc := CORS([]string{"*"})
	corsFunc(c)

	assert.Equal(t, "http://any.domain.com", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORS_OPTIONS(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("OPTIONS", "/", nil)
	c.Request.Header.Set("Origin", "http://localhost:3000")

	corsFunc := CORS([]string{"http://localhost:3000"})
	corsFunc(c)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.True(t, c.IsAborted())
}

func TestCORS_Headers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	corsFunc := CORS([]string{"*"})
	corsFunc(c)

	assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Headers"))
	assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Methods"))
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "POST")
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "GET")
}

func TestCORS_MultipleOrigins(t *testing.T) {
	gin.SetMode(gin.TestMode)

	origins := []string{
		"http://localhost:3000",
		"http://localhost:8080",
		"https://example.com",
	}

	for _, origin := range origins {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Request.Header.Set("Origin", origin)

		corsFunc := CORS(origins)
		corsFunc(c)

		assert.Equal(t, origin, w.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestCORS_NoOriginHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	corsFunc := CORS([]string{"http://localhost:3000"})
	corsFunc(c)

	assert.Equal(t, "", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
}
