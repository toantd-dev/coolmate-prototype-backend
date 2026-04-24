package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestOrderHandler_GetCart(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewOrderHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/cart", nil)

	handler.GetCart(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOrderHandler_AddToCart(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewOrderHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/cart", nil)

	handler.AddToCart(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOrderHandler_UpdateCartItem(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewOrderHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PUT", "/cart/items/1", nil)

	handler.UpdateCartItem(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOrderHandler_RemoveFromCart(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewOrderHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("DELETE", "/cart/items/1", nil)

	handler.RemoveFromCart(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOrderHandler_ApplyCoupon(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewOrderHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/cart/coupon", nil)

	handler.ApplyCoupon(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOrderHandler_Checkout(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewOrderHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/checkout", nil)

	handler.Checkout(c)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestOrderHandler_ListCustomerOrders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewOrderHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/orders", nil)

	handler.ListCustomerOrders(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOrderHandler_GetOrderDetail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewOrderHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/orders/1", nil)

	handler.GetOrderDetail(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOrderHandler_ListVendorOrders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewOrderHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/vendor/orders", nil)

	handler.ListVendorOrders(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOrderHandler_UpdateOrderStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewOrderHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PUT", "/orders/1/status", nil)

	handler.UpdateOrderStatus(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOrderHandler_ListAllOrders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewOrderHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/admin/orders", nil)

	handler.ListAllOrders(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOrderHandler_InitiateReturn(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewOrderHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/returns", nil)

	handler.InitiateReturn(c)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestOrderHandler_ListReturns(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewOrderHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/returns", nil)

	handler.ListReturns(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOrderHandler_NewOrderHandler(t *testing.T) {
	handler := NewOrderHandler(nil)
	assert.NotNil(t, handler)
}
