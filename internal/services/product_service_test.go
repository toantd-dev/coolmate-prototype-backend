package services

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/coolmate/ecommerce-backend/internal/models"
	"github.com/coolmate/ecommerce-backend/pkg/cache"
)

type MockProductRepo struct {
	mock.Mock
}

func (m *MockProductRepo) GetCategory(id uint) (*models.Category, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *MockProductRepo) GetProduct(id uint) (*models.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductRepo) CreateProduct(product *models.Product) (*models.Product, error) {
	args := m.Called(product)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductRepo) UpdateProduct(product *models.Product) (*models.Product, error) {
	args := m.Called(product)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductRepo) DeleteProduct(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockProductRepo) ListProducts(page, limit int, filters map[string]interface{}) ([]*models.Product, int64, error) {
	args := m.Called(page, limit, filters)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*models.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductRepo) ListVendorProducts(vendorID uint, page, limit int) ([]*models.Product, int64, error) {
	args := m.Called(vendorID, page, limit)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*models.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductRepo) GetProductBySlug(slug string) (*models.Product, error) {
	args := m.Called(slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductRepo) GetByID(id uint) (*models.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductRepo) Create(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepo) Update(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepo) List(status string, categoryID uint, limit, offset int) ([]models.Product, int64, error) {
	args := m.Called(status, categoryID, limit, offset)
	if args.Get(1) == nil {
		return args.Get(0).([]models.Product), 0, args.Error(2)
	}
	return args.Get(0).([]models.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductRepo) ListByVendor(vendorID uint, limit, offset int) ([]models.Product, int64, error) {
	args := m.Called(vendorID, limit, offset)
	if args.Get(1) == nil {
		return args.Get(0).([]models.Product), 0, args.Error(2)
	}
	return args.Get(0).([]models.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductRepo) UpdateStatus(id uint, status string) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *MockProductRepo) ListCategories() ([]models.Category, error) {
	args := m.Called()
	return args.Get(0).([]models.Category), args.Error(1)
}

func (m *MockProductRepo) ListPendingApproval(limit, offset int) ([]models.Product, int64, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]models.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductRepo) GetBySlug(slug string) (*models.Product, error) {
	args := m.Called(slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductRepo) CreateImage(image *models.ProductImage) error {
	args := m.Called(image)
	return args.Error(0)
}

func (m *MockProductRepo) GetImages(productID uint) ([]models.ProductImage, error) {
	args := m.Called(productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.ProductImage), args.Error(1)
}

func (m *MockProductRepo) CreateVariant(variant *models.ProductVariant) error {
	args := m.Called(variant)
	return args.Error(0)
}

func (m *MockProductRepo) GetVariants(productID uint) ([]models.ProductVariant, error) {
	args := m.Called(productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.ProductVariant), args.Error(1)
}

func (m *MockProductRepo) SearchProducts(query string, limit int, offset int) ([]models.Product, int64, error) {
	args := m.Called(query, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]models.Product), args.Get(1).(int64), args.Error(2)
}

// Test GetProductByID

func TestGetProductByID_FromCache(t *testing.T) {
	mockProductRepo := new(MockProductRepo)
	mockVendorRepo := new(MockVendorRepository)
	cacheManager := cache.NewCacheManager(nil)

	service := NewProductService(mockProductRepo, mockVendorRepo, cacheManager)

	product := &models.Product{ID: 1, Name: "Test Product"}
	mockProductRepo.On("GetByID", uint(1)).Return(product, nil)

	result, err := service.GetProductByID(1)

	require.NoError(t, err)
	assert.Equal(t, product, result)
	mockProductRepo.AssertExpectations(t)
}

func TestGetProductByID_FromDatabase(t *testing.T) {
	mockProductRepo := new(MockProductRepo)
	mockVendorRepo := new(MockVendorRepository)
	cacheManager := cache.NewCacheManager(nil)

	service := NewProductService(mockProductRepo, mockVendorRepo, cacheManager)

	product := &models.Product{ID: 1, Name: "Test Product"}
	mockProductRepo.On("GetByID", uint(1)).Return(product, nil)

	result, err := service.GetProductByID(1)

	require.NoError(t, err)
	assert.Equal(t, product, result)
	mockProductRepo.AssertExpectations(t)
}

func TestGetProductByID_NotFound(t *testing.T) {
	mockProductRepo := new(MockProductRepo)
	mockVendorRepo := new(MockVendorRepository)
	cacheManager := cache.NewCacheManager(nil)

	service := NewProductService(mockProductRepo, mockVendorRepo, cacheManager)

	mockProductRepo.On("GetByID", uint(999)).Return(nil, errors.New("not found"))

	result, err := service.GetProductByID(999)

	require.Error(t, err)
	assert.Nil(t, result)
}

// Test ValidateProduct

func TestValidateProduct_Valid(t *testing.T) {
	mockProductRepo := new(MockProductRepo)
	mockVendorRepo := new(MockVendorRepository)
	cacheManager := cache.NewCacheManager(nil)

	service := NewProductService(mockProductRepo, mockVendorRepo, cacheManager)

	costPrice := 50.0
	product := &models.Product{
		Name:       "Valid Product",
		VendorID:   1,
		CategoryID: 1,
		BasePrice:  100.0,
		CostPrice:  &costPrice,
	}

	err := service.ValidateProduct(product)
	assert.NoError(t, err)
}

func TestValidateProduct_NilProduct(t *testing.T) {
	mockProductRepo := new(MockProductRepo)
	mockVendorRepo := new(MockVendorRepository)
	cacheManager := cache.NewCacheManager(nil)

	service := NewProductService(mockProductRepo, mockVendorRepo, cacheManager)

	err := service.ValidateProduct(nil)
	assert.Error(t, err)
}

func TestValidateProduct_InvalidName(t *testing.T) {
	mockProductRepo := new(MockProductRepo)
	mockVendorRepo := new(MockVendorRepository)
	cacheManager := cache.NewCacheManager(nil)

	service := NewProductService(mockProductRepo, mockVendorRepo, cacheManager)

	tests := []struct {
		name string
		desc string
	}{
		{"", "empty name"},
		{"AB", "too short"},
		{"A" + string(make([]byte, 300)), "too long"},
	}

	for _, tt := range tests {
		costPrice := 50.0
		product := &models.Product{
			Name:       tt.name,
			VendorID:   1,
			CategoryID: 1,
			BasePrice:  100.0,
			CostPrice:  &costPrice,
		}

		err := service.ValidateProduct(product)
		assert.Error(t, err, tt.desc)
	}
}

func TestValidateProduct_MissingVendorID(t *testing.T) {
	mockProductRepo := new(MockProductRepo)
	mockVendorRepo := new(MockVendorRepository)
	cacheManager := cache.NewCacheManager(nil)

	service := NewProductService(mockProductRepo, mockVendorRepo, cacheManager)

	product := &models.Product{
		Name:       "Test",
		VendorID:   0,
		CategoryID: 1,
		BasePrice:  100.0,
	}

	err := service.ValidateProduct(product)
	assert.Error(t, err)
}

func TestValidateProduct_MissingCategoryID(t *testing.T) {
	mockProductRepo := new(MockProductRepo)
	mockVendorRepo := new(MockVendorRepository)
	cacheManager := cache.NewCacheManager(nil)

	service := NewProductService(mockProductRepo, mockVendorRepo, cacheManager)

	product := &models.Product{
		Name:       "Test",
		VendorID:   1,
		CategoryID: 0,
		BasePrice:  100.0,
	}

	err := service.ValidateProduct(product)
	assert.Error(t, err)
}

func TestValidateProduct_InvalidBasePrice(t *testing.T) {
	mockProductRepo := new(MockProductRepo)
	mockVendorRepo := new(MockVendorRepository)
	cacheManager := cache.NewCacheManager(nil)

	service := NewProductService(mockProductRepo, mockVendorRepo, cacheManager)

	tests := []float64{0, -10}

	for _, price := range tests {
		product := &models.Product{
			Name:       "Test",
			VendorID:   1,
			CategoryID: 1,
			BasePrice:  price,
		}

		err := service.ValidateProduct(product)
		assert.Error(t, err)
	}
}

func TestValidateProduct_NegativeCostPrice(t *testing.T) {
	mockProductRepo := new(MockProductRepo)
	mockVendorRepo := new(MockVendorRepository)
	cacheManager := cache.NewCacheManager(nil)

	service := NewProductService(mockProductRepo, mockVendorRepo, cacheManager)

	costPrice := -10.0
	product := &models.Product{
		Name:       "Test",
		VendorID:   1,
		CategoryID: 1,
		BasePrice:  100.0,
		CostPrice:  &costPrice,
	}

	err := service.ValidateProduct(product)
	assert.Error(t, err)
}

func TestValidateProduct_CostPriceGreaterThanBasePrice(t *testing.T) {
	mockProductRepo := new(MockProductRepo)
	mockVendorRepo := new(MockVendorRepository)
	cacheManager := cache.NewCacheManager(nil)

	service := NewProductService(mockProductRepo, mockVendorRepo, cacheManager)

	costPrice := 150.0
	product := &models.Product{
		Name:       "Test",
		VendorID:   1,
		CategoryID: 1,
		BasePrice:  100.0,
		CostPrice:  &costPrice,
	}

	err := service.ValidateProduct(product)
	assert.Error(t, err)
}

func TestValidateProduct_ReturnableWithoutWindow(t *testing.T) {
	mockProductRepo := new(MockProductRepo)
	mockVendorRepo := new(MockVendorRepository)
	cacheManager := cache.NewCacheManager(nil)

	service := NewProductService(mockProductRepo, mockVendorRepo, cacheManager)

	costPrice := 50.0
	product := &models.Product{
		Name:             "Test",
		VendorID:         1,
		CategoryID:       1,
		BasePrice:        100.0,
		CostPrice:        &costPrice,
		IsReturnable:     true,
		ReturnWindowDays: 0,
	}

	err := service.ValidateProduct(product)
	assert.Error(t, err)
}

// Test CreateProduct

func TestCreateProduct_Success(t *testing.T) {
	mockProductRepo := new(MockProductRepo)
	mockVendorRepo := new(MockVendorRepository)
	cacheManager := cache.NewCacheManager(nil)

	service := NewProductService(mockProductRepo, mockVendorRepo, cacheManager)

	costPrice := 50.0
	product := &models.Product{
		Name:       "Test Product",
		VendorID:   1,
		CategoryID: 1,
		BasePrice:  100.0,
		CostPrice:  &costPrice,
	}

	mockProductRepo.On("Create", mock.MatchedBy(func(p *models.Product) bool {
		return p.Name == "Test Product" && p.Status == "draft"
	})).Return(nil)

	err := service.CreateProduct(product)

	require.NoError(t, err)
	assert.Equal(t, "draft", product.Status)
	mockProductRepo.AssertExpectations(t)
}

func TestCreateProduct_ValidationFailed(t *testing.T) {
	mockProductRepo := new(MockProductRepo)
	mockVendorRepo := new(MockVendorRepository)
	cacheManager := cache.NewCacheManager(nil)

	service := NewProductService(mockProductRepo, mockVendorRepo, cacheManager)

	product := &models.Product{
		Name:       "",
		VendorID:   1,
		CategoryID: 1,
		BasePrice:  100.0,
	}

	err := service.CreateProduct(product)
	assert.Error(t, err)
	mockProductRepo.AssertNotCalled(t, "Create")
}

// Test UpdateProduct

func TestUpdateProduct_PublishedToApproval(t *testing.T) {
	mockProductRepo := new(MockProductRepo)
	mockVendorRepo := new(MockVendorRepository)
	cacheManager := cache.NewCacheManager(nil)

	service := NewProductService(mockProductRepo, mockVendorRepo, cacheManager)

	costPrice := 50.0
	product := &models.Product{
		ID:         1,
		Name:       "Updated Product",
		VendorID:   1,
		CategoryID: 1,
		BasePrice:  100.0,
		CostPrice:  &costPrice,
		Status:     "published",
	}

	mockProductRepo.On("Update", mock.MatchedBy(func(p *models.Product) bool {
		return p.Status == "pending_approval"
	})).Return(nil)

	err := service.UpdateProduct(product)

	require.NoError(t, err)
	assert.Equal(t, "pending_approval", product.Status)
	mockProductRepo.AssertExpectations(t)
}

func TestUpdateProduct_ValidationFailed(t *testing.T) {
	mockProductRepo := new(MockProductRepo)
	mockVendorRepo := new(MockVendorRepository)
	cacheManager := cache.NewCacheManager(nil)

	service := NewProductService(mockProductRepo, mockVendorRepo, cacheManager)

	product := &models.Product{
		ID:       1,
		Name:     "",
		VendorID: 1,
	}

	err := service.UpdateProduct(product)
	assert.Error(t, err)
	mockProductRepo.AssertNotCalled(t, "Update")
}

// Test ApproveProduct

func TestApproveProduct_Success(t *testing.T) {
	mockProductRepo := new(MockProductRepo)
	mockVendorRepo := new(MockVendorRepository)
	cacheManager := cache.NewCacheManager(nil)

	service := NewProductService(mockProductRepo, mockVendorRepo, cacheManager)

	mockProductRepo.On("UpdateStatus", uint(1), "published").Return(nil)

	err := service.ApproveProduct(1)

	require.NoError(t, err)
	mockProductRepo.AssertExpectations(t)
}

// Test RejectProduct

func TestRejectProduct_Success(t *testing.T) {
	mockProductRepo := new(MockProductRepo)
	mockVendorRepo := new(MockVendorRepository)
	cacheManager := cache.NewCacheManager(nil)

	service := NewProductService(mockProductRepo, mockVendorRepo, cacheManager)

	mockProductRepo.On("UpdateStatus", uint(1), "rejected").Return(nil)

	err := service.RejectProduct(1, "Invalid images")

	require.NoError(t, err)
	mockProductRepo.AssertExpectations(t)
}

// Test GetCategories

func TestGetCategories_FromCache(t *testing.T) {
	mockProductRepo := new(MockProductRepo)
	mockVendorRepo := new(MockVendorRepository)
	cacheManager := cache.NewCacheManager(nil)

	service := NewProductService(mockProductRepo, mockVendorRepo, cacheManager)

	categories := []models.Category{
		{ID: 1, Name: "Electronics"},
		{ID: 2, Name: "Fashion"},
	}

	mockProductRepo.On("ListCategories").Return(categories, nil)

	result, err := service.GetCategories()

	require.NoError(t, err)
	assert.Equal(t, categories, result)
	mockProductRepo.AssertExpectations(t)
}

func TestGetCategories_FromDatabase(t *testing.T) {
	mockProductRepo := new(MockProductRepo)
	mockVendorRepo := new(MockVendorRepository)
	cacheManager := cache.NewCacheManager(nil)

	service := NewProductService(mockProductRepo, mockVendorRepo, cacheManager)

	categories := []models.Category{
		{ID: 1, Name: "Electronics"},
	}

	mockProductRepo.On("ListCategories").Return(categories, nil)

	result, err := service.GetCategories()

	require.NoError(t, err)
	assert.Equal(t, categories, result)
	mockProductRepo.AssertExpectations(t)
}

// Test ListProducts

func TestListProducts_Success(t *testing.T) {
	mockProductRepo := new(MockProductRepo)
	mockVendorRepo := new(MockVendorRepository)
	cacheManager := cache.NewCacheManager(nil)

	service := NewProductService(mockProductRepo, mockVendorRepo, cacheManager)

	products := []models.Product{
		{ID: 1, Name: "Product 1"},
	}

	mockProductRepo.On("List", "active", uint(0), 10, 0).Return(products, int64(1), nil)

	result, total, err := service.ListProducts("", 10, 0)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), total)
	mockProductRepo.AssertExpectations(t)
}

// Test ListVendorProducts

func TestListVendorProducts_Success(t *testing.T) {
	mockProductRepo := new(MockProductRepo)
	mockVendorRepo := new(MockVendorRepository)
	cacheManager := cache.NewCacheManager(nil)

	service := NewProductService(mockProductRepo, mockVendorRepo, cacheManager)

	products := []models.Product{
		{ID: 1, Name: "Product 1", VendorID: 1},
	}

	mockProductRepo.On("ListByVendor", uint(1), 10, 0).Return(products, int64(1), nil)

	result, total, err := service.ListVendorProducts(1, 10, 0)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), total)
	mockProductRepo.AssertExpectations(t)
}
