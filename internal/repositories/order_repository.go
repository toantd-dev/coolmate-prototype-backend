package repositories

import (
	"github.com/coolmate/ecommerce-backend/internal/models"
	"gorm.io/gorm"
)

type IOrderRepository interface {
	CreateOrder(order *models.Order) error
	GetOrderByID(id uint) (*models.Order, error)
	ListOrdersByCustomer(customerID uint, limit int, offset int) ([]models.Order, int64, error)
	ListAllOrders(limit int, offset int) ([]models.Order, int64, error)
	UpdateOrderStatus(orderID uint, status string) error
	CreateSubOrder(subOrder *models.SubOrder) error
	ListSubOrdersByVendor(vendorID uint, limit int, offset int) ([]models.SubOrder, int64, error)
	UpdateSubOrderStatus(subOrderID uint, status string) error
	CreateOrderItem(item *models.OrderItem) error
	CreateCart(cart *models.Cart) error
	GetCart(cartID uint) (*models.Cart, error)
	GetCartByUser(userID uint) (*models.Cart, error)
	AddCartItem(item *models.CartItem) error
	RemoveCartItem(itemID uint) error
	UpdateCartItem(item *models.CartItem) error
	ClearCart(cartID uint) error
	CreateReturnRequest(req *models.ReturnRequest) error
	GetReturnRequest(id uint) (*models.ReturnRequest, error)
	ListReturns(customerID uint) ([]models.ReturnRequest, error)
	UpdateReturnStatus(id uint, status string) error
	CreateRefund(refund *models.Refund) error
}

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) IOrderRepository {
	return &OrderRepository{db: db}
}

func (or *OrderRepository) CreateOrder(order *models.Order) error {
	return or.db.Create(order).Error
}

func (or *OrderRepository) GetOrderByID(id uint) (*models.Order, error) {
	var order models.Order
	if err := or.db.Preload("Items").Preload("SubOrders").Preload("StatusLogs").First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (or *OrderRepository) ListOrdersByCustomer(customerID uint, limit int, offset int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	query := or.db.Where("customer_id = ?", customerID)

	if err := query.Model(&models.Order{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Items").Limit(limit).Offset(offset).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (or *OrderRepository) ListAllOrders(limit int, offset int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	if err := or.db.Model(&models.Order{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := or.db.Preload("Items").Preload("SubOrders").Limit(limit).Offset(offset).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (or *OrderRepository) UpdateOrderStatus(orderID uint, status string) error {
	return or.db.Model(&models.Order{}).Where("id = ?", orderID).Update("status", status).Error
}

func (or *OrderRepository) CreateSubOrder(subOrder *models.SubOrder) error {
	return or.db.Create(subOrder).Error
}

func (or *OrderRepository) ListSubOrdersByVendor(vendorID uint, limit int, offset int) ([]models.SubOrder, int64, error) {
	var subOrders []models.SubOrder
	var total int64

	query := or.db.Where("vendor_id = ?", vendorID)

	if err := query.Model(&models.SubOrder{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Order").Limit(limit).Offset(offset).Find(&subOrders).Error; err != nil {
		return nil, 0, err
	}

	return subOrders, total, nil
}

func (or *OrderRepository) UpdateSubOrderStatus(subOrderID uint, status string) error {
	return or.db.Model(&models.SubOrder{}).Where("id = ?", subOrderID).Update("status", status).Error
}

func (or *OrderRepository) CreateOrderItem(item *models.OrderItem) error {
	return or.db.Create(item).Error
}

func (or *OrderRepository) CreateCart(cart *models.Cart) error {
	return or.db.Create(cart).Error
}

func (or *OrderRepository) GetCart(cartID uint) (*models.Cart, error) {
	var cart models.Cart
	if err := or.db.Preload("Items").Preload("Items.Product").Preload("Items.Variant").First(&cart, cartID).Error; err != nil {
		return nil, err
	}
	return &cart, nil
}

func (or *OrderRepository) GetCartByUser(userID uint) (*models.Cart, error) {
	var cart models.Cart
	if err := or.db.Where("user_id = ?", userID).Preload("Items").Preload("Items.Product").Preload("Items.Variant").First(&cart).Error; err != nil {
		return nil, err
	}
	return &cart, nil
}

func (or *OrderRepository) AddCartItem(item *models.CartItem) error {
	return or.db.Create(item).Error
}

func (or *OrderRepository) RemoveCartItem(itemID uint) error {
	return or.db.Delete(&models.CartItem{}, itemID).Error
}

func (or *OrderRepository) UpdateCartItem(item *models.CartItem) error {
	return or.db.Save(item).Error
}

func (or *OrderRepository) ClearCart(cartID uint) error {
	return or.db.Where("cart_id = ?", cartID).Delete(&models.CartItem{}).Error
}

func (or *OrderRepository) CreateReturnRequest(req *models.ReturnRequest) error {
	return or.db.Create(req).Error
}

func (or *OrderRepository) GetReturnRequest(id uint) (*models.ReturnRequest, error) {
	var req models.ReturnRequest
	if err := or.db.First(&req, id).Error; err != nil {
		return nil, err
	}
	return &req, nil
}

func (or *OrderRepository) ListReturns(customerID uint) ([]models.ReturnRequest, error) {
	var returns []models.ReturnRequest
	if err := or.db.Where("customer_id = ?", customerID).Find(&returns).Error; err != nil {
		return nil, err
	}
	return returns, nil
}

func (or *OrderRepository) UpdateReturnStatus(id uint, status string) error {
	return or.db.Model(&models.ReturnRequest{}).Where("id = ?", id).Update("status", status).Error
}

func (or *OrderRepository) CreateRefund(refund *models.Refund) error {
	return or.db.Create(refund).Error
}
