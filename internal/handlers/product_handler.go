package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/coolmate/ecommerce-backend/internal/services"
	"github.com/coolmate/ecommerce-backend/internal/utils"
)

type ProductHandler struct {
	productService services.IProductService
}

func NewProductHandler(productService interface{}) *ProductHandler {
	ps, ok := productService.(services.IProductService)
	if !ok {
		ps = nil
	}
	return &ProductHandler{
		productService: ps,
	}
}

func (ph *ProductHandler) CreateProduct(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusCreated, "Create product", nil)
}

func (ph *ProductHandler) ListVendorProducts(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "List vendor products", nil)
}

func (ph *ProductHandler) UpdateProduct(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Update product", nil)
}

func (ph *ProductHandler) ArchiveProduct(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Archive product", nil)
}

func (ph *ProductHandler) BulkImportProducts(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Bulk import products", nil)
}

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
	products, total, err := ph.productService.ListProducts(search, perPage, offset)
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

func (ph *ProductHandler) ListPendingProducts(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "List pending products", nil)
}

func (ph *ProductHandler) ApproveProduct(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Approve product", nil)
}

func (ph *ProductHandler) RejectProduct(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Reject product", nil)
}

func (ph *ProductHandler) GetCategories(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Get categories", nil)
}
