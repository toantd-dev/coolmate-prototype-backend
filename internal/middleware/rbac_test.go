package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/coolmate/ecommerce-backend/internal/models"
)

func TestRequireRole_Authorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Set("role", string(models.RoleAdmin))

	middlewareFunc := RequireRole(models.RoleAdmin)
	middlewareFunc(c)

	assert.False(t, c.IsAborted())
}

func TestRequireRole_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Set("role", string(models.RoleCustomer))

	middlewareFunc := RequireRole(models.RoleAdmin)
	middlewareFunc(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.True(t, c.IsAborted())
}

func TestRequireRole_NoAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	middlewareFunc := RequireRole(models.RoleAdmin)
	middlewareFunc(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.True(t, c.IsAborted())
}

func TestRequireRole_MultipleAllowedRoles(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Set("role", string(models.RoleVendor))

	middlewareFunc := RequireRole(models.RoleAdmin, models.RoleVendor)
	middlewareFunc(c)

	assert.False(t, c.IsAborted())
}

func TestRequireAdmin_WithAdminRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Set("role", string(models.RoleAdmin))

	middlewareFunc := RequireAdmin()
	middlewareFunc(c)

	assert.False(t, c.IsAborted())
}

func TestRequireAdmin_WithSuperAdminRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Set("role", string(models.RoleSuperAdmin))

	middlewareFunc := RequireAdmin()
	middlewareFunc(c)

	assert.False(t, c.IsAborted())
}

func TestRequireAdmin_WithCustomerRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Set("role", string(models.RoleCustomer))

	middlewareFunc := RequireAdmin()
	middlewareFunc(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.True(t, c.IsAborted())
}

func TestRequireVendor(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Set("role", string(models.RoleVendor))

	middlewareFunc := RequireVendor()
	middlewareFunc(c)

	assert.False(t, c.IsAborted())
}

func TestRequireVendor_WithCustomerRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Set("role", string(models.RoleCustomer))

	middlewareFunc := RequireVendor()
	middlewareFunc(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.True(t, c.IsAborted())
}

func TestRequireSuperAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Set("role", string(models.RoleSuperAdmin))

	middlewareFunc := RequireSuperAdmin()
	middlewareFunc(c)

	assert.False(t, c.IsAborted())
}

func TestRequireSuperAdmin_WithAdminRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Set("role", string(models.RoleAdmin))

	middlewareFunc := RequireSuperAdmin()
	middlewareFunc(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.True(t, c.IsAborted())
}
