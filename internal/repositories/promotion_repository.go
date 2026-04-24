package repositories

import (
	"github.com/coolmate/ecommerce-backend/internal/models"
	"gorm.io/gorm"
)

type IPromotionRepository interface {
	Create(promotion *models.Promotion) error
	GetByID(id uint) (*models.Promotion, error)
	GetByCode(code string) (*models.Promotion, error)
	ListActive(limit int, offset int) ([]models.Promotion, int64, error)
	ListByVendor(vendorID uint) ([]models.Promotion, error)
	Update(promotion *models.Promotion) error
	Delete(id uint) error
	IncrementUsedCount(id uint) error
}

type PromotionRepository struct {
	db *gorm.DB
}

func NewPromotionRepository(db *gorm.DB) IPromotionRepository {
	return &PromotionRepository{db: db}
}

func (pr *PromotionRepository) Create(promotion *models.Promotion) error {
	return pr.db.Create(promotion).Error
}

func (pr *PromotionRepository) GetByID(id uint) (*models.Promotion, error) {
	var promotion models.Promotion
	if err := pr.db.First(&promotion, id).Error; err != nil {
		return nil, err
	}
	return &promotion, nil
}

func (pr *PromotionRepository) GetByCode(code string) (*models.Promotion, error) {
	var promotion models.Promotion
	if err := pr.db.Where("code = ? AND is_active = ?", code, true).First(&promotion).Error; err != nil {
		return nil, err
	}
	return &promotion, nil
}

func (pr *PromotionRepository) ListActive(limit int, offset int) ([]models.Promotion, int64, error) {
	var promotions []models.Promotion
	var total int64

	query := pr.db.Where("is_active = ?", true)

	if err := query.Model(&models.Promotion{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Find(&promotions).Error; err != nil {
		return nil, 0, err
	}

	return promotions, total, nil
}

func (pr *PromotionRepository) ListByVendor(vendorID uint) ([]models.Promotion, error) {
	var promotions []models.Promotion
	if err := pr.db.Where("vendor_id = ?", vendorID).Find(&promotions).Error; err != nil {
		return nil, err
	}
	return promotions, nil
}

func (pr *PromotionRepository) Update(promotion *models.Promotion) error {
	return pr.db.Save(promotion).Error
}

func (pr *PromotionRepository) Delete(id uint) error {
	return pr.db.Delete(&models.Promotion{}, id).Error
}

func (pr *PromotionRepository) IncrementUsedCount(id uint) error {
	return pr.db.Model(&models.Promotion{}).Where("id = ?", id).Update("used_count", gorm.Expr("used_count + ?", 1)).Error
}
