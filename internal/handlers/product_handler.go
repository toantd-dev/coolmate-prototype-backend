package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/coolmate/ecommerce-backend/internal/services"
	"github.com/coolmate/ecommerce-backend/internal/utils"
)

type ProductHandler struct {
	productService *services.ProductService
}

func NewProductHandler(productService interface{}) *ProductHandler {
	ps, ok := productService.(*services.ProductService)
	if !ok {
		ps = nil
	}
	return &ProductHandler{
		productService: ps,
	}
}

// CreateProduct godoc
// @Summary      Create product (STUB)
// @Tags         Products (Vendor)
// @Security     BearerAuth
// @Produce      json
// @Success      201  {object}  utils.APIResponse
// @Router       /vendor/products [post]
func (ph *ProductHandler) CreateProduct(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusCreated, "Create product", nil)
}

// ListVendorProducts godoc
// @Summary      List vendor products (STUB)
// @Tags         Products (Vendor)
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  utils.APIResponse
// @Router       /vendor/products [get]
func (ph *ProductHandler) ListVendorProducts(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "List vendor products", nil)
}

// UpdateProduct godoc
// @Summary      Update product (STUB)
// @Tags         Products (Vendor)
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  utils.APIResponse
// @Router       /vendor/products/{id} [put]
func (ph *ProductHandler) UpdateProduct(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Update product", nil)
}

// ArchiveProduct godoc
// @Summary      Archive product (STUB)
// @Tags         Products (Vendor)
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  utils.APIResponse
// @Router       /vendor/products/{id} [delete]
func (ph *ProductHandler) ArchiveProduct(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Archive product", nil)
}

// BulkImportProducts godoc
// @Summary      Bulk import products from CSV (STUB)
// @Tags         Products (Vendor)
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Success      200  {object}  utils.APIResponse
// @Router       /vendor/products/bulk-import [post]
func (ph *ProductHandler) BulkImportProducts(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Bulk import products", nil)
}

// ListProducts godoc
// @Summary      List published products
// @Description  Public product catalog with pagination and text search.
// @Tags         Products (Public)
// @Produce      json
// @Param        page      query     int     false  "Page number (default 1)"
// @Param        per_page  query     int     false  "Items per page (max 100, default 10)"
// @Param        search    query     string  false  "Substring search on name/description"
// @Success      200       {object}  utils.APIResponse
// @Router       /products [get]
func (ph *ProductHandler) ListProducts(c *gin.Context) {
	// Get query parameters
	pageStr := c.DefaultQuery("page", "1")
	perPageStr := c.DefaultQuery("per_page", "10")
	search := c.Query("search")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	perPage, err := strconv.Atoi(perPageStr)
	if err != nil || perPage < 1 || perPage > 100 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	if ph.productService == nil {
		utils.InternalServerError(c, "Product service not initialized")
		return
	}

	// Get products from service
	products, total, err := ph.productService.ListProducts(search, offset, perPage)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	totalPages := (total + int64(perPage) - 1) / int64(perPage)

	meta := utils.PaginationMeta{
		Total:      total,
		Page:       page,
		PageSize:   perPage,
		TotalPages: totalPages,
	}

	utils.SuccessPaginatedResponse(c, http.StatusOK, "Products retrieved successfully", products, meta)
}

// GetProductBySlug godoc
// @Summary      Get product by slug
// @Description  Fetches a single published product by its URL slug.
// @Tags         Products (Public)
// @Produce      json
// @Param        slug  path      string  true  "Product slug"
// @Success      200   {object}  utils.APIResponse
// @Failure      404   {object}  utils.APIResponse
// @Router       /products/{slug} [get]
func (ph *ProductHandler) GetProductBySlug(c *gin.Context) {
	slug := c.Param("slug")

	if slug == "" {
		utils.BadRequest(c, "Product slug is required")
		return
	}

	if ph.productService == nil {
		utils.InternalServerError(c, "Product service not initialized")
		return
	}

	product, err := ph.productService.GetProductBySlug(slug)
	if err != nil {
		utils.NotFound(c, "Product not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Product retrieved successfully", product)
}

// ListPendingProducts godoc
// @Summary      List products pending approval (STUB)
// @Tags         Products (Admin)
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  utils.APIResponse
// @Router       /admin/products/pending-approval [get]
func (ph *ProductHandler) ListPendingProducts(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "List pending products", nil)
}

// ApproveProduct godoc
// @Summary      Approve product (STUB)
// @Tags         Products (Admin)
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  utils.APIResponse
// @Router       /admin/products/{id}/approve [post]
func (ph *ProductHandler) ApproveProduct(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Approve product", nil)
}

// RejectProduct godoc
// @Summary      Reject product (STUB)
// @Tags         Products (Admin)
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  utils.APIResponse
// @Router       /admin/products/{id}/reject [post]
func (ph *ProductHandler) RejectProduct(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Reject product", nil)
}

// GetCategories godoc
// @Summary      List categories (STUB)
// @Description  Currently a stub — returns no data. Phase 1 TODO.
// @Tags         Products (Public)
// @Produce      json
// @Success      200  {object}  utils.APIResponse
// @Router       /categories [get]
func (ph *ProductHandler) GetCategories(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Get categories", nil)
}
