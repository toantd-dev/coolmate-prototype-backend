package services

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/coolmate/ecommerce-backend/internal/models"
	"github.com/coolmate/ecommerce-backend/pkg/cache"
)

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) CreateOrder(order *models.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderRepository) GetOrderByID(id uint) (*models.Order, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderRepository) ListOrdersByCustomer(customerID uint, limit int, offset int) ([]models.Order, int64, error) {
	args := m.Called(customerID, limit, offset)
	return args.Get(0).([]models.Order), args.Get(1).(int64), args.Error(2)
}

func (m *MockOrderRepository) ListAllOrders(limit int, offset int) ([]models.Order, int64, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]models.Order), args.Get(1).(int64), args.Error(2)
}

func (m *MockOrderRepository) UpdateOrderStatus(orderID uint, status string) error {
	args := m.Called(orderID, status)
	return args.Error(0)
}

func (m *MockOrderRepository) CreateSubOrder(subOrder *models.SubOrder) error {
	args := m.Called(subOrder)
	return args.Error(0)
}

func (m *MockOrderRepository) ListSubOrdersByVendor(vendorID uint, limit int, offset int) ([]models.SubOrder, int64, error) {
	args := m.Called(vendorID, limit, offset)
	return args.Get(0).([]models.SubOrder), args.Get(1).(int64), args.Error(2)
}

func (m *MockOrderRepository) UpdateSubOrderStatus(subOrderID uint, status string) error {
	args := m.Called(subOrderID, status)
	return args.Error(0)
}

func (m *MockOrderRepository) CreateOrderItem(item *models.OrderItem) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MockOrderRepository) CreateCart(cart *models.Cart) error {
	args := m.Called(cart)
	return args.Error(0)
}

func (m *MockOrderRepository) GetCart(cartID uint) (*models.Cart, error) {
	args := m.Called(cartID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Cart), args.Error(1)
}

func (m *MockOrderRepository) GetCartByUser(userID uint) (*models.Cart, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Cart), args.Error(1)
}

func (m *MockOrderRepository) AddCartItem(item *models.CartItem) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MockOrderRepository) RemoveCartItem(itemID uint) error {
	args := m.Called(itemID)
	return args.Error(0)
}

func (m *MockOrderRepository) UpdateCartItem(item *models.CartItem) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MockOrderRepository) ClearCart(cartID uint) error {
	args := m.Called(cartID)
	return args.Error(0)
}

func (m *MockOrderRepository) CreateReturnRequest(req *models.ReturnRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockOrderRepository) GetReturnRequest(id uint) (*models.ReturnRequest, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ReturnRequest), args.Error(1)
}

func (m *MockOrderRepository) ListReturns(customerID uint) ([]models.ReturnRequest, error) {
	args := m.Called(customerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.ReturnRequest), args.Error(1)
}

func (m *MockOrderRepository) UpdateReturnStatus(id uint, status string) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *MockOrderRepository) CreateRefund(refund *models.Refund) error {
	args := m.Called(refund)
	return args.Error(0)
}

type MockPromotionRepository struct {
	mock.Mock
}

func (m *MockPromotionRepository) Create(promotion *models.Promotion) error {
	args := m.Called(promotion)
	return args.Error(0)
}

func (m *MockPromotionRepository) GetByID(id uint) (*models.Promotion, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Promotion), args.Error(1)
}

func (m *MockPromotionRepository) GetByCode(code string) (*models.Promotion, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Promotion), args.Error(1)
}

func (m *MockPromotionRepository) ListActive(limit int, offset int) ([]models.Promotion, int64, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]models.Promotion), args.Get(1).(int64), args.Error(2)
}

func (m *MockPromotionRepository) ListByVendor(vendorID uint) ([]models.Promotion, error) {
	args := m.Called(vendorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Promotion), args.Error(1)
}

func (m *MockPromotionRepository) Update(promotion *models.Promotion) error {
	args := m.Called(promotion)
	return args.Error(0)
}

func (m *MockPromotionRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPromotionRepository) IncrementUsedCount(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockVendorRepository struct {
	mock.Mock
}

func (m *MockVendorRepository) Create(vendor *models.Vendor) error {
	args := m.Called(vendor)
	return args.Error(0)
}

func (m *MockVendorRepository) GetByID(id uint) (*models.Vendor, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Vendor), args.Error(1)
}

func (m *MockVendorRepository) GetByUserID(userID uint) (*models.Vendor, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Vendor), args.Error(1)
}

func (m *MockVendorRepository) GetBySlug(slug string) (*models.Vendor, error) {
	args := m.Called(slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Vendor), args.Error(1)
}

func (m *MockVendorRepository) List(status string, limit int, offset int) ([]models.Vendor, int64, error) {
	args := m.Called(status, limit, offset)
	if args.Get(1) == nil {
		return args.Get(0).([]models.Vendor), 0, args.Error(2)
	}
	return args.Get(0).([]models.Vendor), args.Get(1).(int64), args.Error(2)
}

func (m *MockVendorRepository) Update(vendor *models.Vendor) error {
	args := m.Called(vendor)
	return args.Error(0)
}

func (m *MockVendorRepository) UpdateStatus(vendorID uint, status string) error {
	args := m.Called(vendorID, status)
	return args.Error(0)
}

func (m *MockVendorRepository) CreateDocument(doc *models.VendorDocument) error {
	args := m.Called(doc)
	return args.Error(0)
}

func (m *MockVendorRepository) GetDocuments(vendorID uint) ([]models.VendorDocument, error) {
	args := m.Called(vendorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.VendorDocument), args.Error(1)
}

func (m *MockVendorRepository) CreateBankDetails(details *models.VendorBankDetails) error {
	args := m.Called(details)
	return args.Error(0)
}

func (m *MockVendorRepository) UpdateBankDetails(details *models.VendorBankDetails) error {
	args := m.Called(details)
	return args.Error(0)
}

func (m *MockVendorRepository) GetBankDetails(vendorID uint) (*models.VendorBankDetails, error) {
	args := m.Called(vendorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.VendorBankDetails), args.Error(1)
}

func (m *MockVendorRepository) GetWallet(vendorID uint) (*models.VendorWallet, error) {
	args := m.Called(vendorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.VendorWallet), args.Error(1)
}

func (m *MockVendorRepository) CreateWallet(wallet *models.VendorWallet) error {
	args := m.Called(wallet)
	return args.Error(0)
}

func (m *MockVendorRepository) UpdateWalletBalance(vendorID uint, amount float64) error {
	args := m.Called(vendorID, amount)
	return args.Error(0)
}

func (m *MockVendorRepository) SuspendVendor(vendorID uint) error {
	args := m.Called(vendorID)
	return args.Error(0)
}

// Test SplitOrder

func TestSplitOrder_SingleVendor(t *testing.T) {
	orderRepo := new(MockOrderRepository)
	productRepo := new(MockProductRepository)
	promotionRepo := new(MockPromotionRepository)
	vendorRepo := new(MockVendorRepository)
	cacheManager := cache.NewCacheManager(nil)

	service := NewOrderService(orderRepo, productRepo, promotionRepo, vendorRepo, cacheManager)

	// Mock vendor
	vendorRepo.On("GetByID", uint(1)).Return(&models.Vendor{
		ID:                 1,
		CommissionModel:    "margin",
		CommissionRate:     0.05,
	}, nil)

	productRepo.On("GetCategory", uint(1)).Return(nil, errors.New("not found"))

	order := &models.Order{
		ID: 1,
	}

	cartItems := []models.CartItem{
		{
			ID:        1,
			UnitPrice: 100.0,
			Quantity:  2,
			Product: &models.Product{
				ID:       1,
				VendorID: 1,
			},
		},
	}

	subOrders, err := service.SplitOrder(order, cartItems)

	require.NoError(t, err)
	require.Len(t, subOrders, 1)
	assert.Equal(t, uint(1), subOrders[0].VendorID)
	assert.Equal(t, 200.0, subOrders[0].Subtotal) // 100 * 2
	assert.Equal(t, 10.0, subOrders[0].CommissionAmount) // 200 * 0.05
	assert.Equal(t, 190.0, subOrders[0].VendorEarning) // 200 - 10
}

func TestSplitOrder_MultipleVendors(t *testing.T) {
	orderRepo := new(MockOrderRepository)
	productRepo := new(MockProductRepository)
	promotionRepo := new(MockPromotionRepository)
	vendorRepo := new(MockVendorRepository)
	cacheManager := cache.NewCacheManager(nil)

	service := NewOrderService(orderRepo, productRepo, promotionRepo, vendorRepo, cacheManager)

	vendorRepo.On("GetByID", uint(1)).Return(&models.Vendor{
		ID:              1,
		CommissionModel: "margin",
		CommissionRate:  0.05,
	}, nil)

	vendorRepo.On("GetByID", uint(2)).Return(&models.Vendor{
		ID:              2,
		CommissionModel: "markup",
		CommissionRate:  0.5,
	}, nil)

	productRepo.On("GetCategory", mock.Anything).Return(nil, errors.New("not found"))

	order := &models.Order{
		ID: 1,
	}

	cartItems := []models.CartItem{
		{
			ID:        1,
			UnitPrice: 100.0,
			Quantity:  1,
			Product: &models.Product{
				ID:       1,
				VendorID: 1,
			},
		},
		{
			ID:        2,
			UnitPrice: 50.0,
			Quantity:  2,
			Product: &models.Product{
				ID:       2,
				VendorID: 2,
			},
		},
	}

	subOrders, err := service.SplitOrder(order, cartItems)

	require.NoError(t, err)
	require.Len(t, subOrders, 2)

	// Sort by vendor ID for consistent checking
	vendorMap := make(map[uint]models.SubOrder)
	for _, order := range subOrders {
		vendorMap[order.VendorID] = order
	}

	// Check vendor 1
	vendor1Order := vendorMap[1]
	assert.Equal(t, uint(1), vendor1Order.VendorID)
	assert.Equal(t, 100.0, vendor1Order.Subtotal)

	// Check vendor 2
	vendor2Order := vendorMap[2]
	assert.Equal(t, uint(2), vendor2Order.VendorID)
	assert.Equal(t, 100.0, vendor2Order.Subtotal) // 50 * 2

	// Order totals should be updated
	assert.Equal(t, 200.0, order.Subtotal) // 100 + 100
}

func TestSplitOrder_NilOrder(t *testing.T) {
	orderRepo := new(MockOrderRepository)
	productRepo := new(MockProductRepository)
	promotionRepo := new(MockPromotionRepository)
	vendorRepo := new(MockVendorRepository)
	cacheManager := cache.NewCacheManager(nil)

	service := NewOrderService(orderRepo, productRepo, promotionRepo, vendorRepo, cacheManager)

	subOrders, err := service.SplitOrder(nil, []models.CartItem{})

	require.Error(t, err)
	assert.Nil(t, subOrders)
}

func TestSplitOrder_EmptyCartItems(t *testing.T) {
	orderRepo := new(MockOrderRepository)
	productRepo := new(MockProductRepository)
	promotionRepo := new(MockPromotionRepository)
	vendorRepo := new(MockVendorRepository)
	cacheManager := cache.NewCacheManager(nil)

	service := NewOrderService(orderRepo, productRepo, promotionRepo, vendorRepo, cacheManager)

	order := &models.Order{ID: 1}
	subOrders, err := service.SplitOrder(order, []models.CartItem{})

	require.Error(t, err)
	assert.Nil(t, subOrders)
}

func TestSplitOrder_CartItemWithoutProduct(t *testing.T) {
	orderRepo := new(MockOrderRepository)
	productRepo := new(MockProductRepository)
	promotionRepo := new(MockPromotionRepository)
	vendorRepo := new(MockVendorRepository)
	cacheManager := cache.NewCacheManager(nil)

	service := NewOrderService(orderRepo, productRepo, promotionRepo, vendorRepo, cacheManager)

	order := &models.Order{ID: 1}
	cartItems := []models.CartItem{
		{
			ID:        1,
			UnitPrice: 100.0,
			Quantity:  1,
			Product:   nil,
		},
	}

	subOrders, err := service.SplitOrder(order, cartItems)

	require.Error(t, err)
	assert.Nil(t, subOrders)
}

func TestSplitOrder_VendorNotFound(t *testing.T) {
	orderRepo := new(MockOrderRepository)
	productRepo := new(MockProductRepository)
	promotionRepo := new(MockPromotionRepository)
	vendorRepo := new(MockVendorRepository)
	cacheManager := cache.NewCacheManager(nil)

	service := NewOrderService(orderRepo, productRepo, promotionRepo, vendorRepo, cacheManager)

	vendorRepo.On("GetByID", uint(999)).Return(nil, errors.New("vendor not found"))

	order := &models.Order{ID: 1}
	cartItems := []models.CartItem{
		{
			ID:        1,
			UnitPrice: 100.0,
			Quantity:  1,
			Product: &models.Product{
				ID:       1,
				VendorID: 999,
			},
		},
	}

	subOrders, err := service.SplitOrder(order, cartItems)

	require.Error(t, err)
	assert.Nil(t, subOrders)
}

// Test ValidateStock

func TestValidateStock_SufficientStock(t *testing.T) {
	orderRepo := new(MockOrderRepository)
	productRepo := new(MockProductRepository)
	promotionRepo := new(MockPromotionRepository)
	vendorRepo := new(MockVendorRepository)
	cacheManager := cache.NewCacheManager(nil)

	service := NewOrderService(orderRepo, productRepo, promotionRepo, vendorRepo, cacheManager)

	cartItems := []models.CartItem{
		{
			ID:       1,
			Quantity: 5,
			Variant: &models.ProductVariant{
				Stock: 10,
			},
		},
	}

	err := service.ValidateStock(cartItems)
	require.NoError(t, err)
}

func TestValidateStock_InsufficientStock(t *testing.T) {
	orderRepo := new(MockOrderRepository)
	productRepo := new(MockProductRepository)
	promotionRepo := new(MockPromotionRepository)
	vendorRepo := new(MockVendorRepository)
	cacheManager := cache.NewCacheManager(nil)

	service := NewOrderService(orderRepo, productRepo, promotionRepo, vendorRepo, cacheManager)

	cartItems := []models.CartItem{
		{
			ID:        1,
			Quantity:  15,
			VariantID: new(uint),
			Variant: &models.ProductVariant{
				Stock: 10,
			},
		},
	}

	err := service.ValidateStock(cartItems)
	require.Error(t, err)
}

// Test ApplyPromotions

func TestApplyPromotions_ValidPromo(t *testing.T) {
	orderRepo := new(MockOrderRepository)
	productRepo := new(MockProductRepository)
	promotionRepo := new(MockPromotionRepository)
	vendorRepo := new(MockVendorRepository)
	cacheManager := cache.NewCacheManager(nil)

	service := NewOrderService(orderRepo, productRepo, promotionRepo, vendorRepo, cacheManager)

	promotionRepo.On("GetByCode", "DISCOUNT10").Return(&models.Promotion{
		Code:           "DISCOUNT10",
		DiscountType:   "percent",
		DiscountValue:  10,
		ValidFrom:      time.Now().Add(-1 * time.Hour),
		ValidTo:        time.Now().Add(1 * time.Hour),
	}, nil)

	order := &models.Order{
		Subtotal: 1000.0,
	}

	discount, err := service.ApplyPromotions(order, []string{"DISCOUNT10"})

	require.NoError(t, err)
	assert.Equal(t, 100.0, discount) // 1000 * 10%
}

func TestApplyPromotions_InvalidCode(t *testing.T) {
	orderRepo := new(MockOrderRepository)
	productRepo := new(MockProductRepository)
	promotionRepo := new(MockPromotionRepository)
	vendorRepo := new(MockVendorRepository)
	cacheManager := cache.NewCacheManager(nil)

	service := NewOrderService(orderRepo, productRepo, promotionRepo, vendorRepo, cacheManager)

	promotionRepo.On("GetByCode", "INVALID").Return(nil, errors.New("not found"))

	order := &models.Order{
		Subtotal: 1000.0,
	}

	discount, err := service.ApplyPromotions(order, []string{"INVALID"})

	require.Error(t, err)
	assert.Equal(t, 0.0, discount)
}

func TestApplyPromotions_FlatDiscount(t *testing.T) {
	orderRepo := new(MockOrderRepository)
	productRepo := new(MockProductRepository)
	promotionRepo := new(MockPromotionRepository)
	vendorRepo := new(MockVendorRepository)
	cacheManager := cache.NewCacheManager(nil)

	service := NewOrderService(orderRepo, productRepo, promotionRepo, vendorRepo, cacheManager)

	promotionRepo.On("GetByCode", "FLAT100").Return(&models.Promotion{
		Code:           "FLAT100",
		DiscountType:   "flat",
		DiscountValue:  100,
		ValidFrom:      time.Now().Add(-1 * time.Hour),
		ValidTo:        time.Now().Add(1 * time.Hour),
	}, nil)

	order := &models.Order{
		Subtotal: 1000.0,
	}

	discount, err := service.ApplyPromotions(order, []string{"FLAT100"})

	require.NoError(t, err)
	assert.Equal(t, 100.0, discount)
}

func TestApplyPromotions_MultiplePromos(t *testing.T) {
	orderRepo := new(MockOrderRepository)
	productRepo := new(MockProductRepository)
	promotionRepo := new(MockPromotionRepository)
	vendorRepo := new(MockVendorRepository)
	cacheManager := cache.NewCacheManager(nil)

	service := NewOrderService(orderRepo, productRepo, promotionRepo, vendorRepo, cacheManager)

	promotionRepo.On("GetByCode", "PROMO1").Return(&models.Promotion{
		Code:          "PROMO1",
		DiscountType:  "percent",
		DiscountValue: 10,
	}, nil)

	promotionRepo.On("GetByCode", "PROMO2").Return(&models.Promotion{
		Code:          "PROMO2",
		DiscountType:  "flat",
		DiscountValue: 50,
	}, nil)

	order := &models.Order{
		Subtotal: 1000.0,
	}

	discount, err := service.ApplyPromotions(order, []string{"PROMO1", "PROMO2"})

	require.NoError(t, err)
	assert.Equal(t, 150.0, discount) // 100 + 50
}

// Test GetOrderByID

func TestGetOrderByID_Success(t *testing.T) {
	orderRepo := new(MockOrderRepository)
	productRepo := new(MockProductRepository)
	promotionRepo := new(MockPromotionRepository)
	vendorRepo := new(MockVendorRepository)
	cacheManager := cache.NewCacheManager(nil)

	service := NewOrderService(orderRepo, productRepo, promotionRepo, vendorRepo, cacheManager)

	orderRepo.On("GetOrderByID", uint(1)).Return(&models.Order{
		ID:     1,
		Status: "pending",
	}, nil)

	order, err := service.GetOrderByID(1)

	require.NoError(t, err)
	assert.Equal(t, uint(1), order.ID)
}

func TestGetOrderByID_NotFound(t *testing.T) {
	orderRepo := new(MockOrderRepository)
	productRepo := new(MockProductRepository)
	promotionRepo := new(MockPromotionRepository)
	vendorRepo := new(MockVendorRepository)
	cacheManager := cache.NewCacheManager(nil)

	service := NewOrderService(orderRepo, productRepo, promotionRepo, vendorRepo, cacheManager)

	orderRepo.On("GetOrderByID", uint(999)).Return(nil, errors.New("not found"))

	order, err := service.GetOrderByID(999)

	require.Error(t, err)
	assert.Nil(t, order)
}
