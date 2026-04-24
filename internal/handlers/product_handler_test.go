package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/coolmate/ecommerce-backend/internal/models"
)

type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) ListProducts(search string, limit, offset int) ([]models.Product, int64, error) {
	args := m.Called(search, limit, offset)
	if args.Get(0) == nil {
		return []models.Product{}, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductService) GetProductByID(id uint) (*models.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductService) ListVendorProducts(vendorID uint, limit, offset int) ([]models.Product, int64, error) {
	args := m.Called(vendorID, limit, offset)
	if args.Get(0) == nil {
		return []models.Product{}, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductService) GetProductBySlug(slug string) (*models.Product, error) {
	args := m.Called(slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductService) CreateProduct(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductService) UpdateProduct(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductService) ApproveProduct(productID uint) error {
	args := m.Called(productID)
	return args.Error(0)
}

func (m *MockProductService) RejectProduct(productID uint, reason string) error {
	args := m.Called(productID, reason)
	return args.Error(0)
}

func (m *MockProductService) GetCategories() ([]models.Category, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return []models.Category{}, args.Error(1)
	}
	return args.Get(0).([]models.Category), args.Error(1)
}

func (m *MockProductService) ListPendingApproval(limit, offset int) ([]models.Product, int64, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return []models.Product{}, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductService) ValidateProduct(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func TestProductHandler_CreateProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockProductService{}
	handler := NewProductHandler(mockSvc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/products", nil)

	handler.CreateProduct(c)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestProductHandler_ListVendorProducts(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockProductService{}
	handler := NewProductHandler(mockSvc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/vendor/products", nil)

	handler.ListVendorProducts(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_UpdateProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockProductService{}
	handler := NewProductHandler(mockSvc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PUT", "/products/1", nil)

	handler.UpdateProduct(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_ArchiveProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockProductService{}
	handler := NewProductHandler(mockSvc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("DELETE", "/products/1", nil)

	handler.ArchiveProduct(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_BulkImportProducts(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockProductService{}
	handler := NewProductHandler(mockSvc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/products/bulk", nil)

	handler.BulkImportProducts(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_ListProducts_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockProductService{}

	products := []models.Product{
		{
			Name:       "Laptop",
			BasePrice:  1000.0,
			VendorID:   1,
			CategoryID: 1,
		},
	}
	mockSvc.On("ListProducts", "", 10, 0).Return(products, int64(1), nil)

	handler := NewProductHandler(mockSvc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/products?page=1&per_page=10", nil)

	handler.ListProducts(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestProductHandler_ListProducts_WithSearch(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockProductService{}

	products := []models.Product{
		{Name: "iPhone", BasePrice: 800.0, VendorID: 1, CategoryID: 1},
	}
	mockSvc.On("ListProducts", "iphone", 10, 0).Return(products, int64(1), nil)

	handler := NewProductHandler(mockSvc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/products?search=iphone&page=1&per_page=10", nil)

	handler.ListProducts(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestProductHandler_ListProducts_Page2(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockProductService{}

	products := []models.Product{}
	mockSvc.On("ListProducts", "", 10, 10).Return(products, int64(20), nil)

	handler := NewProductHandler(mockSvc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/products?page=2&per_page=10", nil)

	handler.ListProducts(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestProductHandler_ListProducts_MaxPerPage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockProductService{}

	products := []models.Product{}
	mockSvc.On("ListProducts", "", 10, 0).Return(products, int64(0), nil)

	handler := NewProductHandler(mockSvc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/products?per_page=200", nil)

	handler.ListProducts(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestProductHandler_ListProducts_ServiceNotInitialized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewProductHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/products", nil)

	handler.ListProducts(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestProductHandler_ListProducts_InvalidPage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockProductService{}

	products := []models.Product{}
	mockSvc.On("ListProducts", "", 10, 0).Return(products, int64(0), nil)

	handler := NewProductHandler(mockSvc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/products?page=invalid&per_page=10", nil)

	handler.ListProducts(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestProductHandler_ListProducts_DefaultPagination(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockProductService{}

	products := []models.Product{}
	mockSvc.On("ListProducts", "", 10, 0).Return(products, int64(0), nil)

	handler := NewProductHandler(mockSvc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/products", nil)

	handler.ListProducts(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestProductHandler_ListProducts_MultipleResults(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockProductService{}

	products := []models.Product{
		{ID: 1, Name: "Product 1", BasePrice: 100, VendorID: 1, CategoryID: 1},
		{ID: 2, Name: "Product 2", BasePrice: 200, VendorID: 1, CategoryID: 1},
		{ID: 3, Name: "Product 3", BasePrice: 300, VendorID: 2, CategoryID: 2},
	}
	mockSvc.On("ListProducts", "", 10, 0).Return(products, int64(3), nil)

	handler := NewProductHandler(mockSvc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/products", nil)

	handler.ListProducts(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestProductHandler_NewProductHandler_WithService(t *testing.T) {
	mockSvc := &MockProductService{}
	handler := NewProductHandler(mockSvc)

	assert.NotNil(t, handler)
	assert.Equal(t, mockSvc, handler.productService)
}

func TestProductHandler_NewProductHandler_WithInvalidType(t *testing.T) {
	handler := NewProductHandler("invalid")

	assert.NotNil(t, handler)
	assert.Nil(t, handler.productService)
}

func TestProductHandler_ListProducts_EmptyResults(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockProductService{}

	products := []models.Product{}
	mockSvc.On("ListProducts", "", 10, 0).Return(products, int64(0), nil)

	handler := NewProductHandler(mockSvc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/products", nil)

	handler.ListProducts(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestProductHandler_ListProducts_LargePerPage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockProductService{}

	products := []models.Product{}
	mockSvc.On("ListProducts", "", 10, 0).Return(products, int64(0), nil)

	handler := NewProductHandler(mockSvc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/products?per_page=500", nil)

	handler.ListProducts(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

