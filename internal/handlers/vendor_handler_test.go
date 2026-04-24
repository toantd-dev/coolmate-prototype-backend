package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestVendorHandler_ListVendors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewVendorHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/vendors", nil)

	handler.ListVendors(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestVendorHandler_GetVendor(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewVendorHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/vendors/1", nil)

	handler.GetVendor(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestVendorHandler_GetProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewVendorHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/vendor/profile", nil)

	handler.GetProfile(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestVendorHandler_UpdateProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewVendorHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PUT", "/vendor/profile", nil)

	handler.UpdateProfile(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestVendorHandler_UploadDocument(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewVendorHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/vendor/documents", nil)

	handler.UploadDocument(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestVendorHandler_ApproveVendor(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewVendorHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/vendors/1/approve", nil)

	handler.ApproveVendor(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestVendorHandler_RejectVendor(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewVendorHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/vendors/1/reject", nil)

	handler.RejectVendor(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestVendorHandler_SuspendVendor(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewVendorHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/vendors/1/suspend", nil)

	handler.SuspendVendor(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestVendorHandler_GetStaff(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewVendorHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/vendor/staff", nil)

	handler.GetStaff(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestVendorHandler_CreateStaff(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewVendorHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/vendor/staff", nil)

	handler.CreateStaff(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestVendorHandler_UpdateStaff(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewVendorHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PUT", "/vendor/staff/1", nil)

	handler.UpdateStaff(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestVendorHandler_DeleteStaff(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewVendorHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("DELETE", "/vendor/staff/1", nil)

	handler.DeleteStaff(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestVendorHandler_AcceptAgreement(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewVendorHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/vendor/agreement", nil)

	handler.AcceptAgreement(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestVendorHandler_UpdateBankDetails(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewVendorHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PUT", "/vendor/bank", nil)

	handler.UpdateBankDetails(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestVendorHandler_UpdateCommission(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewVendorHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PUT", "/vendors/1/commission", nil)

	handler.UpdateCommission(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestVendorHandler_GetVendorWallet(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewVendorHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/vendor/wallet", nil)

	handler.GetVendorWallet(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestVendorHandler_SettleVendor(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewVendorHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/vendors/1/settle", nil)

	handler.SettleVendor(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestVendorHandler_GetSettlements(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewVendorHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/vendor/settlements", nil)

	handler.GetSettlements(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestVendorHandler_NewVendorHandler(t *testing.T) {
	handler := NewVendorHandler(nil)
	assert.NotNil(t, handler)
}
