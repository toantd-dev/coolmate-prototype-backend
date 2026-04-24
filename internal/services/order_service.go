package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/coolmate/ecommerce-backend/internal/models"
	"github.com/coolmate/ecommerce-backend/internal/repositories"
	"github.com/coolmate/ecommerce-backend/pkg/cache"
)

type OrderService struct {
	orderRepo       repositories.IOrderRepository
	productRepo     repositories.IProductRepository
	promotionRepo   repositories.IPromotionRepository
	vendorRepo      repositories.IVendorRepository
	cache           *cache.CacheManager
	commissionSvc   *CommissionService
}

func NewOrderService(
	orderRepo repositories.IOrderRepository,
	productRepo repositories.IProductRepository,
	promotionRepo repositories.IPromotionRepository,
	vendorRepo repositories.IVendorRepository,
	cache *cache.CacheManager,
) *OrderService {
	return &OrderService{
		orderRepo:     orderRepo,
		productRepo:   productRepo,
		promotionRepo: promotionRepo,
		vendorRepo:    vendorRepo,
		cache:         cache,
		commissionSvc: NewCommissionService(productRepo),
	}
}

func (os *OrderService) GetOrderByID(id uint) (*models.Order, error) {
	return os.orderRepo.GetOrderByID(id)
}

func (os *OrderService) ListOrders(limit int, offset int) ([]models.Order, int64, error) {
	return os.orderRepo.ListAllOrders(limit, offset)
}

// SplitOrder splits a cart into vendor-specific sub-orders
// This is critical for multi-vendor checkout
func (os *OrderService) SplitOrder(
	order *models.Order,
	cartItems []models.CartItem,
) ([]models.SubOrder, error) {
	if order == nil {
		return nil, errors.New("order cannot be nil")
	}

	if len(cartItems) == 0 {
		return nil, errors.New("cart items cannot be empty")
	}

	// Group items by vendor
	vendorItems := make(map[uint][]models.CartItem)
	for _, item := range cartItems {
		if item.Product == nil {
			return nil, fmt.Errorf("cart item %d has no product", item.ID)
		}
		vendorID := item.Product.VendorID
		vendorItems[vendorID] = append(vendorItems[vendorID], item)
	}

	var subOrders []models.SubOrder

	// Create sub-order for each vendor
	for vendorID, items := range vendorItems {
		var subOrderTotal float64
		var commissionTotal float64
		var vendorEarning float64

		// Calculate totals and commissions for this vendor's items
		for _, item := range items {
			itemTotal := item.UnitPrice * float64(item.Quantity)
			subOrderTotal += itemTotal

			// Calculate commission for this item
			orderItem := &models.OrderItem{
				Product:      item.Product,
				UnitPrice:    item.UnitPrice,
				Quantity:     item.Quantity,
			}

			vendor, err := os.vendorRepo.GetByID(vendorID)
			if err != nil {
				return nil, fmt.Errorf("vendor %d not found: %w", vendorID, err)
			}

			commission, _, err := os.commissionSvc.CalculateCommission(
				orderItem,
				vendor.CommissionModel,
				vendor.CommissionRate,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to calculate commission: %w", err)
			}

			commissionTotal += commission
		}

		vendorEarning = subOrderTotal - commissionTotal

		subOrder := models.SubOrder{
			OrderID:         order.ID,
			VendorID:        vendorID,
			Status:          "pending",
			Subtotal:        subOrderTotal,
			CommissionAmount: commissionTotal,
			VendorEarning:   vendorEarning,
		}

		subOrders = append(subOrders, subOrder)
	}

	// Update order totals
	for _, subOrder := range subOrders {
		order.Subtotal += subOrder.Subtotal
	}
	order.GrandTotal = order.Subtotal - order.DiscountTotal + order.ShippingTotal

	return subOrders, nil
}

// ValidateStock checks if all cart items have sufficient stock
func (os *OrderService) ValidateStock(cartItems []models.CartItem) error {
	for _, item := range cartItems {
		if item.Variant == nil || item.Variant.Stock < item.Quantity {
			variantID := uint(0)
			if item.VariantID != nil {
				variantID = *item.VariantID
			}
			return fmt.Errorf("insufficient stock for product variant %d", variantID)
		}
	}
	return nil
}

// ApplyPromotions applies promotions to an order and returns total discount
func (os *OrderService) ApplyPromotions(
	order *models.Order,
	promoCodes []string,
) (float64, error) {
	ctx := context.Background()
	var totalDiscount float64

	for _, code := range promoCodes {
		// Check cache first
		cacheKey := fmt.Sprintf("promo:%s", code)
		promo := &models.Promotion{}
		err := os.cache.Get(ctx, cacheKey, promo)
		if err == nil && os.isPromotionValid(promo) {
			discount := os.calculateDiscount(order.Subtotal, promo)
			totalDiscount += discount
			continue
		}

		// Fetch from database if not cached
		promo, err = os.promotionRepo.GetByCode(code)
		if err != nil {
			return 0, fmt.Errorf("promotion code %s not found", code)
		}

		if !os.isPromotionValid(promo) {
			return 0, fmt.Errorf("promotion code %s is not valid", code)
		}

		// Cache for 1 hour
		os.cache.Set(ctx, cacheKey, promo, 3600*time.Second)

		discount := os.calculateDiscount(order.Subtotal, promo)
		totalDiscount += discount
	}

	return totalDiscount, nil
}

// isPromotionValid checks if promotion is currently valid
func (os *OrderService) isPromotionValid(promo *models.Promotion) bool {
	if promo == nil {
		return false
	}
	// Check usage limit, expiry date, etc.
	// Placeholder - implement based on your Promotion model
	return true
}

// calculateDiscount calculates discount amount based on promotion
func (os *OrderService) calculateDiscount(subtotal float64, promo *models.Promotion) float64 {
	if promo.DiscountType == "percent" {
		return subtotal * (promo.DiscountValue / 100)
	}
	// Flat discount
	return promo.DiscountValue
}
