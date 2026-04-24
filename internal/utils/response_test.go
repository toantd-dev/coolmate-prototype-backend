package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSuccessResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := map[string]interface{}{"id": 1, "name": "test"}
	SuccessResponse(c, http.StatusOK, "Success", data)

	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.True(t, response.Success)
	assert.Equal(t, "Success", response.Message)
}

func TestSuccessResponse_Created(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	SuccessResponse(c, http.StatusCreated, "Created", map[string]string{"id": "123"})

	assert.Equal(t, http.StatusCreated, w.Code)

	var response APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.True(t, response.Success)
}

func TestSuccessResponse_Nil(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	SuccessResponse(c, http.StatusOK, "OK", nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.True(t, response.Success)
	assert.Nil(t, response.Data)
}

func TestSuccessPaginatedResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	items := []string{"a", "b", "c"}
	meta := PaginationMeta{
		Total:      100,
		Page:       1,
		PageSize:   10,
		TotalPages: 10,
	}

	SuccessPaginatedResponse(c, http.StatusOK, "Items", items, meta)

	assert.Equal(t, http.StatusOK, w.Code)

	var response PaginatedResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.True(t, response.Success)
	assert.Equal(t, int64(100), response.Meta.Total)
	assert.Equal(t, int64(10), response.Meta.TotalPages)
}

func TestErrorResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	errors := map[string]string{"field": "error message"}
	ErrorResponse(c, http.StatusBadRequest, "Validation failed", errors)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.False(t, response.Success)
	assert.Equal(t, "Validation failed", response.Message)
}

func TestBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	BadRequest(c, "Invalid input")

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.False(t, response.Success)
	assert.Equal(t, "Invalid input", response.Message)
}

func TestUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Unauthorized(c, "Invalid credentials")

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.False(t, response.Success)
	assert.Equal(t, "Invalid credentials", response.Message)
}

func TestForbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Forbidden(c, "Access denied")

	assert.Equal(t, http.StatusForbidden, w.Code)

	var response APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.False(t, response.Success)
	assert.Equal(t, "Access denied", response.Message)
}

func TestNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	NotFound(c, "Resource not found")

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.False(t, response.Success)
	assert.Equal(t, "Resource not found", response.Message)
}

func TestInternalServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	InternalServerError(c, "Database error")

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.False(t, response.Success)
	assert.Equal(t, "Database error", response.Message)
}

func TestAPIResponse_Structure(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	testData := map[string]interface{}{
		"id":   1,
		"name": "test",
		"tags": []string{"a", "b"},
	}

	SuccessResponse(c, http.StatusOK, "Test message", testData)

	var response APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.True(t, response.Success)
	assert.Equal(t, "Test message", response.Message)
	assert.NotNil(t, response.Data)
}

func TestPaginationMeta_Structure(t *testing.T) {
	meta := PaginationMeta{
		Total:      100,
		Page:       2,
		PageSize:   20,
		TotalPages: 5,
	}

	assert.Equal(t, int64(100), meta.Total)
	assert.Equal(t, 2, meta.Page)
	assert.Equal(t, 20, meta.PageSize)
	assert.Equal(t, int64(5), meta.TotalPages)
}

func TestErrorResponse_WithNilErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	ErrorResponse(c, http.StatusInternalServerError, "Error occurred", nil)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.False(t, response.Success)
	assert.Nil(t, response.Errors)
}

func TestSuccessResponse_WithComplexData(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	complexData := map[string]interface{}{
		"user": map[string]interface{}{
			"id":    1,
			"email": "test@example.com",
		},
		"settings": map[string]interface{}{
			"notifications": true,
			"theme":         "dark",
		},
	}

	SuccessResponse(c, http.StatusOK, "User data", complexData)

	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.True(t, response.Success)
	assert.NotNil(t, response.Data)
}
