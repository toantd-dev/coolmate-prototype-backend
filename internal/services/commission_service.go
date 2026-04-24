package services

import (
	"errors"
	"fmt"

	"github.com/coolmate/ecommerce-backend/internal/models"
	"github.com/coolmate/ecommerce-backend/internal/repositories"
)

type CommissionService struct {
	productRepo repositories.IProductRepository
}

func NewCommissionService(
	productRepo repositories.IProductRepository,
) *CommissionService {
	return &CommissionService{
		productRepo: productRepo,
	}
}

// CommissionConfig holds configuration for a specific commission calculation
type CommissionConfig struct {
	CommissionModel string
	CommissionRate  float64
}

// CalculateCommission calculates commission for an order item based on priority hierarchy
// Priority: Category commission > Vendor commission > Platform default (5% margin)
func (cs *CommissionService) CalculateCommission(
	orderItem *models.OrderItem,
	vendorCommissionModel string,
	vendorCommissionRate float64,
) (commissionAmount float64, config CommissionConfig, err error) {
	if orderItem == nil {
		return 0, CommissionConfig{}, errors.New("order item cannot be nil")
	}

	// Priority 1: Category commission
	if orderItem.Product != nil && orderItem.Product.CategoryID > 0 {
		category, err := cs.productRepo.GetCategory(orderItem.Product.CategoryID)
		if err == nil && category != nil && category.CommissionModel != "" && category.CommissionRate != nil {
			commissionAmount = cs.calculateByModel(
				orderItem.UnitPrice,
				orderItem.Quantity,
				category.CommissionModel,
				*category.CommissionRate,
			)
			return commissionAmount, CommissionConfig{
				CommissionModel: category.CommissionModel,
				CommissionRate:  *category.CommissionRate,
			}, nil
		}
	}

	// Priority 2: Vendor commission
	if vendorCommissionModel != "" {
		commissionAmount = cs.calculateByModel(
			orderItem.UnitPrice,
			orderItem.Quantity,
			vendorCommissionModel,
			vendorCommissionRate,
		)
		return commissionAmount, CommissionConfig{
			CommissionModel: vendorCommissionModel,
			CommissionRate:  vendorCommissionRate,
		}, nil
	}

	// Priority 3: Platform default (5% margin model)
	platformRate := 0.05
	commissionAmount = cs.calculateByModel(
		orderItem.UnitPrice,
		orderItem.Quantity,
		"margin",
		platformRate,
	)
	return commissionAmount, CommissionConfig{
		CommissionModel: "margin",
		CommissionRate:  platformRate,
	}, nil
}

// calculateByModel calculates commission based on model type
// Margin: commission = subtotal * rate
// Markup: commission = subtotal / (1 + rate) * (1 - (1 / (1 + rate)))
func (cs *CommissionService) calculateByModel(
	unitPrice float64,
	quantity int,
	model string,
	rate float64,
) float64 {
	subtotal := unitPrice * float64(quantity)

	switch model {
	case "margin":
		// Margin model: flat percentage of subtotal
		return subtotal * rate

	case "markup":
		// Markup model: based on markup percentage
		// If vendor marks up 50%, commission = subtotal / 1.5
		return subtotal / (1 + rate) * (1 - (1 / (1 + rate)))

	default:
		// Default to margin model if unknown
		return subtotal * rate
	}
}

// ValidateCommissionRate ensures commission rate is within acceptable range
func (cs *CommissionService) ValidateCommissionRate(rate float64, model string) error {
	if rate < 0 || rate > 1 {
		return fmt.Errorf("commission rate must be between 0 and 1, got %f", rate)
	}

	if model == "markup" && rate < 0 {
		return errors.New("markup rate cannot be negative")
	}

	return nil
}

// GetCategoryCommission retrieves commission configuration for a category
func (cs *CommissionService) GetCategoryCommission(categoryID uint) (CommissionConfig, error) {
	category, err := cs.productRepo.GetCategory(categoryID)
	if err != nil {
		return CommissionConfig{}, err
	}

	rate := 0.0
	if category.CommissionRate != nil {
		rate = *category.CommissionRate
	}

	return CommissionConfig{
		CommissionModel: category.CommissionModel,
		CommissionRate:  rate,
	}, nil
}
