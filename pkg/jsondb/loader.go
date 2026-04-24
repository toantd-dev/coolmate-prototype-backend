package jsondb

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/coolmate/ecommerce-backend/internal/models"
)

type ProductData struct {
	ID          uint                    `json:"id"`
	Name        string                  `json:"name"`
	Slug        string                  `json:"slug"`
	Description string                  `json:"description"`
	BasePrice   float64                 `json:"price"`
	SKU         string                  `json:"sku"`
	CategoryID  uint                    `json:"category_id"`
	BrandID     *uint                   `json:"brand_id"`
	VendorID    uint                    `json:"vendor_id"`
	Status      string                  `json:"status"`
	Vendor      *VendorData             `json:"vendor"`
	Images      []ProductImageData      `json:"images"`
	Variants    []ProductVariantData    `json:"variants"`
	CreatedAt   string                  `json:"created_at"`
	UpdatedAt   string                  `json:"updated_at"`
}

type VendorData struct {
	ID           uint   `json:"id"`
	BusinessName string `json:"business_name"`
	Logo         *string `json:"logo"`
}

type ProductImageData struct {
	ID        uint   `json:"id"`
	ProductID uint   `json:"product_id"`
	URL       string `json:"url"`
	IsPrimary bool   `json:"is_primary"`
	SortOrder int    `json:"sort_order"`
}

type ProductVariantData struct {
	ID        uint    `json:"id"`
	ProductID uint    `json:"product_id"`
	SKU       string  `json:"sku"`
	Price     float64 `json:"price"`
	Stock     int     `json:"stock"`
}

type ProductsFile struct {
	Products []ProductData `json:"products"`
}

type JSONLoader struct {
	filePath string
	data     *ProductsFile
}

func NewJSONLoader(filePath string) (*JSONLoader, error) {
	loader := &JSONLoader{filePath: filePath}
	if err := loader.Load(); err != nil {
		return nil, err
	}
	return loader, nil
}

func (jl *JSONLoader) Load() error {
	// Get absolute path if relative
	path := jl.filePath
	if !filepath.IsAbs(path) {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}
		path = filepath.Join(wd, path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read JSON file %s: %w", path, err)
	}

	var productsFile ProductsFile
	if err := json.Unmarshal(data, &productsFile); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	jl.data = &productsFile
	return nil
}

func (jl *JSONLoader) GetProducts() []ProductData {
	if jl.data == nil {
		return []ProductData{}
	}
	return jl.data.Products
}

func (jl *JSONLoader) GetProductByID(id uint) *ProductData {
	if jl.data == nil {
		return nil
	}
	for _, p := range jl.data.Products {
		if p.ID == id {
			return &p
		}
	}
	return nil
}

func (jl *JSONLoader) GetProductBySlug(slug string) *ProductData {
	if jl.data == nil {
		return nil
	}
	for _, p := range jl.data.Products {
		if p.Slug == slug {
			return &p
		}
	}
	return nil
}

// ToModel converts ProductData to models.Product
func (pd *ProductData) ToModel() *models.Product {
	createdAt, _ := time.Parse(time.RFC3339, pd.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, pd.UpdatedAt)

	product := &models.Product{
		ID:          pd.ID,
		Name:        pd.Name,
		Slug:        pd.Slug,
		Description: pd.Description,
		BasePrice:   pd.BasePrice,
		SKU:         pd.SKU,
		CategoryID:  pd.CategoryID,
		BrandID:     pd.BrandID,
		VendorID:    pd.VendorID,
		Status:      pd.Status,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	// Convert images
	for _, img := range pd.Images {
		product.Images = append(product.Images, models.ProductImage{
			ID:        img.ID,
			ProductID: img.ProductID,
			URL:       img.URL,
			IsPrimary: img.IsPrimary,
			SortOrder: img.SortOrder,
		})
	}

	// Convert variants
	for _, variant := range pd.Variants {
		product.Variants = append(product.Variants, models.ProductVariant{
			ID:        variant.ID,
			ProductID: variant.ProductID,
			SKU:       variant.SKU,
			Price:     variant.Price,
			Stock:     variant.Stock,
		})
	}

	return product
}
