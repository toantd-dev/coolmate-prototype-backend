package services

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/coolmate/ecommerce-backend/internal/models"
)

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) GetByID(id uint) (*models.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductRepository) GetBySlug(slug string) (*models.Product, error) {
	args := m.Called(slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductRepository) ListByVendor(vendorID uint, limit int, offset int) ([]models.Product, int64, error) {
	args := m.Called(vendorID, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]models.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductRepository) List(status string, categoryID uint, limit int, offset int) ([]models.Product, int64, error) {
	args := m.Called(status, categoryID, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]models.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductRepository) Update(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) UpdateStatus(productID uint, status string) error {
	args := m.Called(productID, status)
	return args.Error(0)
}

func (m *MockProductRepository) CreateImage(image *models.ProductImage) error {
	args := m.Called(image)
	return args.Error(0)
}

func (m *MockProductRepository) GetImages(productID uint) ([]models.ProductImage, error) {
	args := m.Called(productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.ProductImage), args.Error(1)
}

func (m *MockProductRepository) CreateVariant(variant *models.ProductVariant) error {
	args := m.Called(variant)
	return args.Error(0)
}

func (m *MockProductRepository) GetVariants(productID uint) ([]models.ProductVariant, error) {
	args := m.Called(productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.ProductVariant), args.Error(1)
}

func (m *MockProductRepository) ListPendingApproval(limit int, offset int) ([]models.Product, int64, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]models.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductRepository) GetCategory(id uint) (*models.Category, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *MockProductRepository) ListCategories() ([]models.Category, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Category), args.Error(1)
}

func (m *MockProductRepository) SearchProducts(query string, limit int, offset int) ([]models.Product, int64, error) {
	args := m.Called(query, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]models.Product), args.Get(1).(int64), args.Error(2)
}

// Test CalculateCommission with priority hierarchy

func TestCalculateCommission_CategoryCommission(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewCommissionService(mockRepo)

	categoryID := uint(1)
	commissionRate := 0.1
	mockRepo.On("GetCategory", categoryID).Return(&models.Category{
		ID:              categoryID,
		CommissionModel: "margin",
		CommissionRate:  &commissionRate,
	}, nil)

	orderItem := &models.OrderItem{
		ID:        1,
		UnitPrice: 100.0,
		Quantity:  2,
		Product: &models.Product{
			CategoryID: categoryID,
		},
	}

	commission, config, err := service.CalculateCommission(
		orderItem,
		"margin",
		0.05,
	)

	require.NoError(t, err)
	assert.Equal(t, 20.0, commission) // 200 * 0.1
	assert.Equal(t, "margin", config.CommissionModel)
	assert.Equal(t, 0.1, config.CommissionRate)
	mockRepo.AssertExpectations(t)
}

func TestCalculateCommission_VendorCommissionWhenNoCategoryCommission(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewCommissionService(mockRepo)

	mockRepo.On("GetCategory", uint(1)).Return(nil, errors.New("not found"))

	orderItem := &models.OrderItem{
		ID:        1,
		UnitPrice: 100.0,
		Quantity:  2,
		Product: &models.Product{
			CategoryID: 1,
		},
	}

	commission, config, err := service.CalculateCommission(
		orderItem,
		"markup",
		0.5,
	)

	require.NoError(t, err)
	// Markup model: 200 / (1 + 0.5) * (1 - (1 / (1 + 0.5))) = 200 / 1.5 * (1 - 2/3) = 133.33 * 0.33 = 44.44
	assert.InDelta(t, 44.44, commission, 0.01)
	assert.Equal(t, "markup", config.CommissionModel)
	assert.Equal(t, 0.5, config.CommissionRate)
	mockRepo.AssertExpectations(t)
}

func TestCalculateCommission_PlatformDefaultWhenNoVendorCommission(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewCommissionService(mockRepo)

	mockRepo.On("GetCategory", uint(1)).Return(nil, errors.New("not found"))

	orderItem := &models.OrderItem{
		ID:        1,
		UnitPrice: 100.0,
		Quantity:  2,
		Product: &models.Product{
			CategoryID: 1,
		},
	}

	commission, config, err := service.CalculateCommission(
		orderItem,
		"", // No vendor commission
		0,
	)

	require.NoError(t, err)
	assert.Equal(t, 10.0, commission) // 200 * 0.05 (default platform margin)
	assert.Equal(t, "margin", config.CommissionModel)
	assert.Equal(t, 0.05, config.CommissionRate)
	mockRepo.AssertExpectations(t)
}

func TestCalculateCommission_NilOrderItem(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewCommissionService(mockRepo)

	commission, config, err := service.CalculateCommission(nil, "margin", 0.05)

	require.Error(t, err)
	assert.Equal(t, 0.0, commission)
	assert.Equal(t, CommissionConfig{}, config)
}

func TestCalculateCommission_NoProduct(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewCommissionService(mockRepo)

	orderItem := &models.OrderItem{
		ID:        1,
		UnitPrice: 100.0,
		Quantity:  2,
		Product:   nil,
	}

	commission, _, err := service.CalculateCommission(
		orderItem,
		"margin",
		0.05,
	)

	require.NoError(t, err)
	assert.Equal(t, 10.0, commission) // Falls back to vendor/platform commission
}

// Test calculateByModel

func TestCalculateByModel_Margin(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewCommissionService(mockRepo)

	commission := service.calculateByModel(100.0, 5, "margin", 0.1)
	assert.Equal(t, 50.0, commission) // 500 * 0.1
}

func TestCalculateByModel_Markup(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewCommissionService(mockRepo)

	// Markup with rate 0.5 means 50% markup
	commission := service.calculateByModel(100.0, 2, "markup", 0.5)
	expected := 200.0 / (1 + 0.5) * (1 - (1 / (1 + 0.5)))
	assert.InDelta(t, expected, commission, 0.01)
}

func TestCalculateByModel_UnknownModel(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewCommissionService(mockRepo)

	commission := service.calculateByModel(100.0, 5, "unknown", 0.1)
	assert.Equal(t, 50.0, commission) // Defaults to margin
}

// Test ValidateCommissionRate

func TestValidateCommissionRate_Valid(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewCommissionService(mockRepo)

	tests := []struct {
		rate  float64
		model string
	}{
		{0, "margin"},
		{0.5, "margin"},
		{1, "margin"},
		{0.5, "markup"},
	}

	for _, tt := range tests {
		err := service.ValidateCommissionRate(tt.rate, tt.model)
		assert.NoError(t, err)
	}
}

func TestValidateCommissionRate_Invalid(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewCommissionService(mockRepo)

	tests := []struct {
		rate  float64
		model string
	}{
		{-0.1, "margin"},
		{1.1, "margin"},
		{-0.5, "markup"},
		{1.5, "markup"},
	}

	for _, tt := range tests {
		err := service.ValidateCommissionRate(tt.rate, tt.model)
		assert.Error(t, err)
	}
}

// Test GetCategoryCommission

func TestGetCategoryCommission_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewCommissionService(mockRepo)

	rate := 0.15
	mockRepo.On("GetCategory", uint(1)).Return(&models.Category{
		ID:              1,
		CommissionModel: "margin",
		CommissionRate:  &rate,
	}, nil)

	config, err := service.GetCategoryCommission(1)

	require.NoError(t, err)
	assert.Equal(t, "margin", config.CommissionModel)
	assert.Equal(t, 0.15, config.CommissionRate)
	mockRepo.AssertExpectations(t)
}

func TestGetCategoryCommission_NotFound(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewCommissionService(mockRepo)

	mockRepo.On("GetCategory", uint(999)).Return(nil, errors.New("category not found"))

	config, err := service.GetCategoryCommission(999)

	assert.Error(t, err)
	assert.Equal(t, CommissionConfig{}, config)
	mockRepo.AssertExpectations(t)
}

// Edge case tests

func TestCalculateCommission_ZeroQuantity(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewCommissionService(mockRepo)

	mockRepo.On("GetCategory", uint(1)).Return(nil, errors.New("not found"))

	orderItem := &models.OrderItem{
		ID:        1,
		UnitPrice: 100.0,
		Quantity:  0,
		Product: &models.Product{
			CategoryID: 1,
		},
	}

	commission, _, err := service.CalculateCommission(
		orderItem,
		"margin",
		0.05,
	)

	require.NoError(t, err)
	assert.Equal(t, 0.0, commission)
}

func TestCalculateCommission_ZeroUnitPrice(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewCommissionService(mockRepo)

	mockRepo.On("GetCategory", uint(1)).Return(nil, errors.New("not found"))

	orderItem := &models.OrderItem{
		ID:        1,
		UnitPrice: 0.0,
		Quantity:  5,
		Product: &models.Product{
			CategoryID: 1,
		},
	}

	commission, _, err := service.CalculateCommission(
		orderItem,
		"margin",
		0.05,
	)

	require.NoError(t, err)
	assert.Equal(t, 0.0, commission)
}

func TestCalculateCommission_LargePrices(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewCommissionService(mockRepo)

	mockRepo.On("GetCategory", uint(1)).Return(nil, errors.New("not found"))

	orderItem := &models.OrderItem{
		ID:        1,
		UnitPrice: 1000000.0,
		Quantity:  100,
		Product: &models.Product{
			CategoryID: 1,
		},
	}

	commission, _, err := service.CalculateCommission(
		orderItem,
		"margin",
		0.05,
	)

	require.NoError(t, err)
	assert.Equal(t, 5000000.0, commission) // 100M * 0.05
}
