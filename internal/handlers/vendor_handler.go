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

// ListVendors godoc
// @Summary      List all vendors (STUB, admin)
// @Tags         Vendors (Admin)
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  utils.APIResponse
// @Router       /admin/vendors [get]
func (vh *VendorHandler) ListVendors(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "List vendors", nil)
}

// GetVendor godoc
// @Summary      Get vendor by ID (STUB, admin)
// @Tags         Vendors (Admin)
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Vendor ID"
// @Success      200  {object}  utils.APIResponse
// @Router       /admin/vendors/{id} [get]
func (vh *VendorHandler) GetVendor(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Get vendor", nil)
}

// GetProfile godoc
// @Summary      Get own vendor profile (STUB)
// @Tags         Vendor
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  utils.APIResponse
// @Router       /vendor/profile [get]
func (vh *VendorHandler) GetProfile(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Get vendor profile", nil)
}

// UpdateProfile godoc
// @Summary      Update own vendor profile (STUB)
// @Tags         Vendor
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  utils.APIResponse
// @Router       /vendor/profile [put]
func (vh *VendorHandler) UpdateProfile(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Update vendor profile", nil)
}

// UploadDocument godoc
// @Summary      Upload KYC document (STUB)
// @Tags         Vendor
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Success      200  {object}  utils.APIResponse
// @Router       /vendor/documents/upload [post]
func (vh *VendorHandler) UploadDocument(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Upload document", nil)
}

// ApproveVendor godoc
// @Summary      Approve vendor (STUB, admin)
// @Tags         Vendors (Admin)
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Vendor ID"
// @Success      200  {object}  utils.APIResponse
// @Router       /admin/vendors/{id}/approve [post]
func (vh *VendorHandler) ApproveVendor(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Approve vendor", nil)
}

// RejectVendor godoc
// @Summary      Reject vendor (STUB, admin)
// @Tags         Vendors (Admin)
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Vendor ID"
// @Success      200  {object}  utils.APIResponse
// @Router       /admin/vendors/{id}/reject [post]
func (vh *VendorHandler) RejectVendor(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Reject vendor", nil)
}

// SuspendVendor godoc
// @Summary      Suspend vendor (STUB, admin)
// @Tags         Vendors (Admin)
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Vendor ID"
// @Success      200  {object}  utils.APIResponse
// @Router       /admin/vendors/{id}/suspend [post]
func (vh *VendorHandler) SuspendVendor(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Suspend vendor", nil)
}

// GetStaff godoc
// @Summary      List vendor staff sub-accounts (STUB)
// @Tags         Vendor
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  utils.APIResponse
// @Router       /vendor/staff [get]
func (vh *VendorHandler) GetStaff(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Get staff", nil)
}

// CreateStaff godoc
// @Summary      Create vendor staff sub-account (STUB)
// @Tags         Vendor
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  utils.APIResponse
// @Router       /vendor/staff [post]
func (vh *VendorHandler) CreateStaff(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Create staff", nil)
}

// UpdateStaff godoc
// @Summary      Update vendor staff sub-account (STUB)
// @Tags         Vendor
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Staff ID"
// @Success      200  {object}  utils.APIResponse
// @Router       /vendor/staff/{id} [put]
func (vh *VendorHandler) UpdateStaff(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Update staff", nil)
}

// DeleteStaff godoc
// @Summary      Delete vendor staff sub-account (STUB)
// @Tags         Vendor
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Staff ID"
// @Success      200  {object}  utils.APIResponse
// @Router       /vendor/staff/{id} [delete]
func (vh *VendorHandler) DeleteStaff(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Delete staff", nil)
}

// AcceptAgreement godoc
// @Summary      Accept current vendor agreement version (STUB)
// @Tags         Vendor
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  utils.APIResponse
// @Router       /vendor/agreement/accept [post]
func (vh *VendorHandler) AcceptAgreement(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Accept agreement", nil)
}

// UpdateBankDetails godoc
// @Summary      Update vendor bank details (STUB, admin only per SRS)
// @Tags         Vendors (Admin)
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Vendor ID"
// @Success      200  {object}  utils.APIResponse
// @Router       /admin/vendors/{id}/bank-details [put]
func (vh *VendorHandler) UpdateBankDetails(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Update bank details", nil)
}

// UpdateCommission godoc
// @Summary      Set per-vendor commission override (STUB, admin)
// @Tags         Vendors (Admin)
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Vendor ID"
// @Success      200  {object}  utils.APIResponse
// @Router       /admin/vendors/{id}/commission [put]
func (vh *VendorHandler) UpdateCommission(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Update commission", nil)
}

// GetVendorWallet godoc
// @Summary      Get vendor wallet (STUB, admin)
// @Tags         Vendors (Admin)
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Vendor ID"
// @Success      200  {object}  utils.APIResponse
// @Router       /admin/vendors/{id}/wallet [get]
func (vh *VendorHandler) GetVendorWallet(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Get vendor wallet", nil)
}

// SettleVendor godoc
// @Summary      Trigger a payout settlement (STUB, admin)
// @Tags         Vendors (Admin)
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Vendor ID"
// @Success      200  {object}  utils.APIResponse
// @Router       /admin/vendors/{id}/settle [post]
func (vh *VendorHandler) SettleVendor(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Settle vendor", nil)
}

// GetSettlements godoc
// @Summary      Get vendor settlement history (STUB, admin)
// @Tags         Vendors (Admin)
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Vendor ID"
// @Success      200  {object}  utils.APIResponse
// @Router       /admin/vendors/{id}/settlements [get]
func (vh *VendorHandler) GetSettlements(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Get settlements", nil)
}
