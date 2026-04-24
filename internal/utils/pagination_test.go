package utils

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetPaginationParams_Defaults(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/?page=1&pageSize=20", nil)

	params := GetPaginationParams(c)

	assert.Equal(t, 1, params.Page)
	assert.Equal(t, 20, params.PageSize)
	assert.Equal(t, "", params.Sort)
}

func TestGetPaginationParams_CustomValues(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/?page=3&pageSize=50&sort=name", nil)

	params := GetPaginationParams(c)

	assert.Equal(t, 3, params.Page)
	assert.Equal(t, 50, params.PageSize)
	assert.Equal(t, "name", params.Sort)
}

func TestGetPaginationParams_InvalidPage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/?page=invalid&pageSize=20", nil)

	params := GetPaginationParams(c)

	assert.Equal(t, 1, params.Page)
	assert.Equal(t, 20, params.PageSize)
}

func TestGetPaginationParams_InvalidLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/?page=1&pageSize=abc", nil)

	params := GetPaginationParams(c)

	assert.Equal(t, 1, params.Page)
	assert.Equal(t, 20, params.PageSize)
}

func TestGetPaginationParams_NegativePage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/?page=-5&pageSize=20", nil)

	params := GetPaginationParams(c)

	assert.Equal(t, 1, params.Page)
}

func TestGetPaginationParams_ZeroPage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/?page=0&pageSize=20", nil)

	params := GetPaginationParams(c)

	assert.Equal(t, 1, params.Page)
}

func TestGetPaginationParams_MaxLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/?page=1&pageSize=500", nil)

	params := GetPaginationParams(c)

	assert.Equal(t, 1, params.Page)
	assert.Equal(t, 100, params.PageSize)
}

func TestGetPaginationParams_NegativeLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/?page=1&pageSize=-10", nil)

	params := GetPaginationParams(c)

	assert.Equal(t, 1, params.Page)
	assert.Equal(t, 20, params.PageSize)
}

func TestGetPaginationParams_EmptyQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	params := GetPaginationParams(c)

	assert.Equal(t, 1, params.Page)
	assert.Equal(t, 20, params.PageSize)
	assert.Equal(t, "", params.Sort)
}

func TestPaginationParams_GetOffset_Page1(t *testing.T) {
	params := PaginationParams{
		Page:     1,
		PageSize: 10,
	}

	offset := params.GetOffset()

	assert.Equal(t, 0, offset)
}

func TestPaginationParams_GetOffset_Page2(t *testing.T) {
	params := PaginationParams{
		Page:     2,
		PageSize: 10,
	}

	offset := params.GetOffset()

	assert.Equal(t, 10, offset)
}

func TestPaginationParams_GetOffset_Page5(t *testing.T) {
	params := PaginationParams{
		Page:     5,
		PageSize: 20,
	}

	offset := params.GetOffset()

	assert.Equal(t, 80, offset)
}

func TestPaginationParams_GetOffset_Page0(t *testing.T) {
	params := PaginationParams{
		Page:     0,
		PageSize: 10,
	}

	offset := params.GetOffset()

	assert.Equal(t, -10, offset)
}

func TestCalculatePaginationMeta_SinglePage(t *testing.T) {
	meta := CalculatePaginationMeta(5, 1, 10)

	assert.Equal(t, int64(5), meta.Total)
	assert.Equal(t, 1, meta.Page)
	assert.Equal(t, 10, meta.PageSize)
	assert.Equal(t, int64(1), meta.TotalPages)
}

func TestCalculatePaginationMeta_MultiplePages(t *testing.T) {
	meta := CalculatePaginationMeta(25, 1, 10)

	assert.Equal(t, int64(25), meta.Total)
	assert.Equal(t, 1, meta.Page)
	assert.Equal(t, 10, meta.PageSize)
	assert.Equal(t, int64(3), meta.TotalPages)
}

func TestCalculatePaginationMeta_ExactMultiple(t *testing.T) {
	meta := CalculatePaginationMeta(100, 2, 20)

	assert.Equal(t, int64(100), meta.Total)
	assert.Equal(t, 2, meta.Page)
	assert.Equal(t, 20, meta.PageSize)
	assert.Equal(t, int64(5), meta.TotalPages)
}

func TestCalculatePaginationMeta_Zero(t *testing.T) {
	meta := CalculatePaginationMeta(0, 1, 10)

	assert.Equal(t, int64(0), meta.Total)
	assert.Equal(t, 1, meta.Page)
	assert.Equal(t, 10, meta.PageSize)
	assert.Equal(t, int64(0), meta.TotalPages)
}

func TestCalculatePaginationMeta_LargeNumbers(t *testing.T) {
	meta := CalculatePaginationMeta(1000000, 50, 100)

	assert.Equal(t, int64(1000000), meta.Total)
	assert.Equal(t, 50, meta.Page)
	assert.Equal(t, 100, meta.PageSize)
	assert.Equal(t, int64(10000), meta.TotalPages)
}

func TestGetPaginationParams_WithSort(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []string{
		"name",
		"created_at",
		"price_desc",
		"rating_asc",
	}

	for _, sortValue := range testCases {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/?sort="+sortValue, nil)

		params := GetPaginationParams(c)

		assert.Equal(t, sortValue, params.Sort)
	}
}
