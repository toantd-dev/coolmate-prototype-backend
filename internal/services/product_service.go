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

type IProductService interface {
	GetProductByID(id uint) (*models.Product, error)
	ListProducts(search string, limit int, offset int) ([]models.Product, int64, error)
	ListVendorProducts(vendorID uint, limit int, offset int) ([]models.Product, int64, error)
	GetProductBySlug(slug string) (*models.Product, error)
	CreateProduct(product *models.Product) error
	UpdateProduct(product *models.Product) error
	ApproveProduct(productID uint) error
	RejectProduct(productID uint, reason string) error
	GetCategories() ([]models.Category, error)
	ListPendingApproval(limit int, offset int) ([]models.Product, int64, error)
	ValidateProduct(product *models.Product) error
}

type ProductService struct {
	productRepo repositories.IProductRepository
	vendorRepo  repositories.IVendorRepository
	cache       *cache.CacheManager
}

func NewProductService(
	productRepo repositories.IProductRepository,
	vendorRepo repositories.IVendorRepository,
	cache *cache.CacheManager,
) *ProductService {
	return &ProductService{
		productRepo: productRepo,
		vendorRepo:  vendorRepo,
		cache:       cache,
	}
}

func (ps *ProductService) GetProductByID(id uint) (*models.Product, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("product:%d", id)

	// Try to get from cache
	product := &models.Product{}
	err := ps.cache.Get(ctx, cacheKey, product)
	if err == nil {
		return product, nil
	}

	// Fetch from database
	product, err = ps.productRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Cache for 5 minutes
	ps.cache.Set(ctx, cacheKey, product, 5*time.Minute)
	return product, nil
}

func (ps *ProductService) ListProducts(search string, limit int, offset int) ([]models.Product, int64, error) {
	if search != "" {
		return ps.productRepo.SearchProducts(search, limit, offset)
	}
	return ps.productRepo.List("active", 0, limit, offset)
}

func (ps *ProductService) ListVendorProducts(vendorID uint, limit int, offset int) ([]models.Product, int64, error) {
	return ps.productRepo.ListByVendor(vendorID, limit, offset)
}

func (ps *ProductService) GetProductBySlug(slug string) (*models.Product, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("product:slug:%s", slug)

	// Try to get from cache
	product := &models.Product{}
	err := ps.cache.Get(ctx, cacheKey, product)
	if err == nil {
		return product, nil
	}

	// Fetch from repository
	product, err = ps.productRepo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}

	// Cache for 5 minutes
	ps.cache.Set(ctx, cacheKey, product, 5*time.Minute)
	return product, nil
}

func (ps *ProductService) CreateProduct(product *models.Product) error {
	ctx := context.Background()
	if err := ps.ValidateProduct(product); err != nil {
		return err
	}

	// New products start as draft
	product.Status = "draft"

	if err := ps.productRepo.Create(product); err != nil {
		return err
	}

	// Invalidate category cache
	ps.cache.Delete(ctx, "categories")

	return nil
}

func (ps *ProductService) UpdateProduct(product *models.Product) error {
	ctx := context.Background()
	if err := ps.ValidateProduct(product); err != nil {
		return err
	}

	// If already published, reset to pending_approval on update
	if product.Status == "published" {
		product.Status = "pending_approval"
	}

	if err := ps.productRepo.Update(product); err != nil {
		return err
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("product:%d", product.ID)
	ps.cache.Delete(ctx, cacheKey, "categories")

	return nil
}

func (ps *ProductService) ApproveProduct(productID uint) error {
	ctx := context.Background()
	if err := ps.productRepo.UpdateStatus(productID, "published"); err != nil {
		return err
	}
	ps.cache.Delete(ctx, fmt.Sprintf("product:%d", productID))
	return nil
}

func (ps *ProductService) RejectProduct(productID uint, reason string) error {
	ctx := context.Background()
	if err := ps.productRepo.UpdateStatus(productID, "rejected"); err != nil {
		return err
	}
	ps.cache.Delete(ctx, fmt.Sprintf("product:%d", productID))
	return nil
}

func (ps *ProductService) GetCategories() ([]models.Category, error) {
	ctx := context.Background()
	// Check cache
	var categories []models.Category
	err := ps.cache.Get(ctx, "categories", &categories)
	if err == nil {
		return categories, nil
	}

	// Fetch from database
	categories, err = ps.productRepo.ListCategories()
	if err != nil {
		return nil, err
	}

	// Cache for 24 hours (categories rarely change)
	ps.cache.Set(ctx, "categories", categories, 24*time.Hour)
	return categories, nil
}

func (ps *ProductService) ListPendingApproval(limit int, offset int) ([]models.Product, int64, error) {
	return ps.productRepo.ListPendingApproval(limit, offset)
}

// ValidateProduct validates product data before creation/update
func (ps *ProductService) ValidateProduct(product *models.Product) error {
	if product == nil {
		return errors.New("product cannot be nil")
	}

	if product.Name == "" || len(product.Name) < 3 || len(product.Name) > 255 {
		return errors.New("product name must be 3-255 characters")
	}

	if product.VendorID == 0 {
		return errors.New("vendor ID is required")
	}

	if product.CategoryID == 0 {
		return errors.New("category ID is required")
	}

	if product.BasePrice <= 0 {
		return errors.New("base price must be greater than 0")
	}

	if product.CostPrice != nil && *product.CostPrice < 0 {
		return errors.New("cost price cannot be negative")
	}

	if product.CostPrice != nil && *product.CostPrice >= product.BasePrice {
		return errors.New("cost price must be less than base price")
	}

	// Validate return window
	if product.IsReturnable && product.ReturnWindowDays <= 0 {
		return errors.New("return window days must be greater than 0")
	}

	return nil
}
