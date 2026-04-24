package repositories

import (
	"github.com/coolmate/ecommerce-backend/internal/models"
	"gorm.io/gorm"
)

type IProductRepository interface {
	Create(product *models.Product) error
	GetByID(id uint) (*models.Product, error)
	GetBySlug(slug string) (*models.Product, error)
	ListByVendor(vendorID uint, limit int, offset int) ([]models.Product, int64, error)
	List(status string, categoryID uint, limit int, offset int) ([]models.Product, int64, error)
	Update(product *models.Product) error
	UpdateStatus(productID uint, status string) error
	CreateImage(image *models.ProductImage) error
	GetImages(productID uint) ([]models.ProductImage, error)
	CreateVariant(variant *models.ProductVariant) error
	GetVariants(productID uint) ([]models.ProductVariant, error)
	ListPendingApproval(limit int, offset int) ([]models.Product, int64, error)
	GetCategory(id uint) (*models.Category, error)
	ListCategories() ([]models.Category, error)
	SearchProducts(query string, limit int, offset int) ([]models.Product, int64, error)
}

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) IProductRepository {
	return &ProductRepository{db: db}
}

func (pr *ProductRepository) Create(product *models.Product) error {
	return pr.db.Create(product).Error
}

func (pr *ProductRepository) GetByID(id uint) (*models.Product, error) {
	var product models.Product
	if err := pr.db.Preload("Vendor").Preload("Category").Preload("Brand").Preload("Images").Preload("Variants").First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (pr *ProductRepository) GetBySlug(slug string) (*models.Product, error) {
	var product models.Product
	if err := pr.db.Preload("Vendor").Preload("Category").Preload("Brand").Preload("Images").Preload("Variants").Where("slug = ? AND status = ?", slug, "published").First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (pr *ProductRepository) ListByVendor(vendorID uint, limit int, offset int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := pr.db.Where("vendor_id = ?", vendorID)

	if err := query.Model(&models.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Category").Preload("Images").Preload("Variants").Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (pr *ProductRepository) List(status string, categoryID uint, limit int, offset int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := pr.db
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	if err := query.Model(&models.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Vendor").Preload("Category").Preload("Images").Preload("Variants").Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (pr *ProductRepository) Update(product *models.Product) error {
	return pr.db.Save(product).Error
}

func (pr *ProductRepository) UpdateStatus(productID uint, status string) error {
	return pr.db.Model(&models.Product{}).Where("id = ?", productID).Update("status", status).Error
}

func (pr *ProductRepository) CreateImage(image *models.ProductImage) error {
	return pr.db.Create(image).Error
}

func (pr *ProductRepository) GetImages(productID uint) ([]models.ProductImage, error) {
	var images []models.ProductImage
	if err := pr.db.Where("product_id = ?", productID).Order("sort_order").Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}

func (pr *ProductRepository) CreateVariant(variant *models.ProductVariant) error {
	return pr.db.Create(variant).Error
}

func (pr *ProductRepository) GetVariants(productID uint) ([]models.ProductVariant, error) {
	var variants []models.ProductVariant
	if err := pr.db.Where("product_id = ?", productID).Find(&variants).Error; err != nil {
		return nil, err
	}
	return variants, nil
}

func (pr *ProductRepository) ListPendingApproval(limit int, offset int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := pr.db.Where("status = ?", "pending_approval")

	if err := query.Model(&models.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Vendor").Preload("Category").Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (pr *ProductRepository) GetCategory(id uint) (*models.Category, error) {
	var category models.Category
	if err := pr.db.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (pr *ProductRepository) ListCategories() ([]models.Category, error) {
	var categories []models.Category
	if err := pr.db.Where("parent_id IS NULL").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (pr *ProductRepository) SearchProducts(query string, limit int, offset int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	searchQuery := "%" + query + "%"
	dbQuery := pr.db.Where("status = ? AND (name ILIKE ? OR description ILIKE ? OR sku ILIKE ?)", "published", searchQuery, searchQuery, searchQuery)

	if err := dbQuery.Model(&models.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := dbQuery.Preload("Vendor").Preload("Category").Preload("Images").Preload("Variants").Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
