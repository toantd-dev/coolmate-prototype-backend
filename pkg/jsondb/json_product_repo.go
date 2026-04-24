package jsondb

import (
	"fmt"
	"strings"

	"github.com/coolmate/ecommerce-backend/internal/models"
)

type JSONProductRepository struct {
	loader *JSONLoader
}

func NewJSONProductRepository(loader *JSONLoader) *JSONProductRepository {
	return &JSONProductRepository{loader: loader}
}

func (jpr *JSONProductRepository) Create(product *models.Product) error {
	return fmt.Errorf("create operation not supported in JSON mode")
}

func (jpr *JSONProductRepository) GetByID(id uint) (*models.Product, error) {
	productData := jpr.loader.GetProductByID(id)
	if productData == nil {
		return nil, fmt.Errorf("product not found")
	}
	return productData.ToModel(), nil
}

func (jpr *JSONProductRepository) GetBySlug(slug string) (*models.Product, error) {
	productData := jpr.loader.GetProductBySlug(slug)
	if productData == nil {
		return nil, fmt.Errorf("product not found")
	}
	return productData.ToModel(), nil
}

func (jpr *JSONProductRepository) ListByVendor(vendorID uint, limit int, offset int) ([]models.Product, int64, error) {
	products := jpr.loader.GetProducts()
	var filtered []ProductData

	for _, p := range products {
		if p.VendorID == vendorID && (p.Status == "active" || p.Status == "draft") {
			filtered = append(filtered, p)
		}
	}

	total := int64(len(filtered))
	var result []models.Product

	// Apply pagination
	if offset < len(filtered) {
		end := offset + limit
		if end > len(filtered) {
			end = len(filtered)
		}
		for _, p := range filtered[offset:end] {
			result = append(result, *p.ToModel())
		}
	}

	return result, total, nil
}

func (jpr *JSONProductRepository) List(status string, categoryID uint, limit int, offset int) ([]models.Product, int64, error) {
	products := jpr.loader.GetProducts()
	var filtered []ProductData

	for _, p := range products {
		// Filter by status (only show active/published products)
		if p.Status != "active" {
			continue
		}

		// Filter by category if specified
		if categoryID > 0 && p.CategoryID != categoryID {
			continue
		}

		filtered = append(filtered, p)
	}

	total := int64(len(filtered))
	var result []models.Product

	// Apply pagination
	if offset < len(filtered) {
		end := offset + limit
		if end > len(filtered) {
			end = len(filtered)
		}
		for _, p := range filtered[offset:end] {
			result = append(result, *p.ToModel())
		}
	}

	return result, total, nil
}

func (jpr *JSONProductRepository) ListByVendorAndStatus(vendorID uint, status string, limit int, offset int) ([]models.Product, int64, error) {
	products := jpr.loader.GetProducts()
	var filtered []ProductData

	for _, p := range products {
		if p.VendorID == vendorID {
			if status == "" || p.Status == status {
				filtered = append(filtered, p)
			}
		}
	}

	total := int64(len(filtered))
	var result []models.Product

	// Apply pagination
	if offset < len(filtered) {
		end := offset + limit
		if end > len(filtered) {
			end = len(filtered)
		}
		for _, p := range filtered[offset:end] {
			result = append(result, *p.ToModel())
		}
	}

	return result, total, nil
}

func (jpr *JSONProductRepository) Update(product *models.Product) error {
	return fmt.Errorf("update operation not supported in JSON mode")
}

func (jpr *JSONProductRepository) UpdateStatus(productID uint, status string) error {
	return fmt.Errorf("update operation not supported in JSON mode")
}

func (jpr *JSONProductRepository) GetByIDs(ids []uint) ([]models.Product, error) {
	products := jpr.loader.GetProducts()
	var result []models.Product

	for _, id := range ids {
		for _, p := range products {
			if p.ID == id {
				result = append(result, *p.ToModel())
				break
			}
		}
	}

	return result, nil
}

func (jpr *JSONProductRepository) SearchProducts(query string, limit int, offset int) ([]models.Product, int64, error) {
	products := jpr.loader.GetProducts()
	var filtered []ProductData

	queryLower := strings.ToLower(query)

	for _, p := range products {
		if p.Status != "active" {
			continue
		}

		// Search in name and description
		if strings.Contains(strings.ToLower(p.Name), queryLower) ||
			strings.Contains(strings.ToLower(p.Description), queryLower) ||
			strings.Contains(strings.ToLower(p.SKU), queryLower) {
			filtered = append(filtered, p)
		}
	}

	total := int64(len(filtered))
	var result []models.Product

	// Apply pagination
	if offset < len(filtered) {
		end := offset + limit
		if end > len(filtered) {
			end = len(filtered)
		}
		for _, p := range filtered[offset:end] {
			result = append(result, *p.ToModel())
		}
	}

	return result, total, nil
}

func (jpr *JSONProductRepository) ListPendingApproval(limit int, offset int) ([]models.Product, int64, error) {
	products := jpr.loader.GetProducts()
	var filtered []ProductData

	for _, p := range products {
		if p.Status == "draft" {
			filtered = append(filtered, p)
		}
	}

	total := int64(len(filtered))
	var result []models.Product

	// Apply pagination
	if offset < len(filtered) {
		end := offset + limit
		if end > len(filtered) {
			end = len(filtered)
		}
		for _, p := range filtered[offset:end] {
			result = append(result, *p.ToModel())
		}
	}

	return result, total, nil
}

func (jpr *JSONProductRepository) CreateImage(image *models.ProductImage) error {
	return fmt.Errorf("create image operation not supported in JSON mode")
}

func (jpr *JSONProductRepository) GetImages(productID uint) ([]models.ProductImage, error) {
	productData := jpr.loader.GetProductByID(productID)
	if productData == nil {
		return []models.ProductImage{}, nil
	}
	return productData.ToModel().Images, nil
}

func (jpr *JSONProductRepository) CreateVariant(variant *models.ProductVariant) error {
	return fmt.Errorf("create variant operation not supported in JSON mode")
}

func (jpr *JSONProductRepository) GetVariants(productID uint) ([]models.ProductVariant, error) {
	productData := jpr.loader.GetProductByID(productID)
	if productData == nil {
		return []models.ProductVariant{}, nil
	}
	return productData.ToModel().Variants, nil
}

func (jpr *JSONProductRepository) GetCategory(id uint) (*models.Category, error) {
	return nil, fmt.Errorf("get category operation not supported in JSON mode")
}

func (jpr *JSONProductRepository) ListCategories() ([]models.Category, error) {
	return []models.Category{}, nil
}
