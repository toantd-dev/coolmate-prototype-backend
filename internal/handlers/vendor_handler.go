package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/coolmate/ecommerce-backend/internal/utils"
)

type VendorHandler struct {
	// vendorService *services.VendorService
}

func NewVendorHandler(vendorService interface{}) *VendorHandler {
	return &VendorHandler{}
}

func (vh *VendorHandler) ListVendors(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "List vendors", nil)
}

func (vh *VendorHandler) GetVendor(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Get vendor", nil)
}

func (vh *VendorHandler) GetProfile(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Get vendor profile", nil)
}

func (vh *VendorHandler) UpdateProfile(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Update vendor profile", nil)
}

func (vh *VendorHandler) UploadDocument(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Upload document", nil)
}

func (vh *VendorHandler) ApproveVendor(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Approve vendor", nil)
}

func (vh *VendorHandler) RejectVendor(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Reject vendor", nil)
}

func (vh *VendorHandler) SuspendVendor(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Suspend vendor", nil)
}

func (vh *VendorHandler) GetStaff(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Get staff", nil)
}

func (vh *VendorHandler) CreateStaff(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Create staff", nil)
}

func (vh *VendorHandler) UpdateStaff(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Update staff", nil)
}

func (vh *VendorHandler) DeleteStaff(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Delete staff", nil)
}

func (vh *VendorHandler) AcceptAgreement(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Accept agreement", nil)
}

func (vh *VendorHandler) UpdateBankDetails(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Update bank details", nil)
}

func (vh *VendorHandler) UpdateCommission(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Update commission", nil)
}

func (vh *VendorHandler) GetVendorWallet(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Get vendor wallet", nil)
}

func (vh *VendorHandler) SettleVendor(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Settle vendor", nil)
}

func (vh *VendorHandler) GetSettlements(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Get settlements", nil)
}
