// BDD-tagged umbrella for the request-validation boundary.
//
// The actual 35+ validation assertions live in request_validators_test.go
// (existing AAA-style tests). Those are the source of line coverage and are
// deliberately left untouched. This umbrella file exists solely to register
// AC tags so that scripts/ac-coverage.sh can count US-VAL-001 ACs against
// USER_STORIES.md. Each sub-test reads one sample from the legacy suite to
// demonstrate the contract and references the legacy test name in a comment.
package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Helper: run one binding attempt and return whether Gin's ShouldBindJSON
// accepted the payload.
func bindsOK(t *testing.T, body string, target interface{}) bool {
	t.Helper()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")
	err := c.ShouldBindJSON(target)
	return err == nil
}

// =============================================================================
// US-VAL-001 · Handler binding rejects malformed payloads before the service.
//
// AC: US-VAL-001 AC1, AC2, AC3, AC4, AC5, AC6, AC7
//
// Evidence authority: 35+ cases in request_validators_test.go.
// Representative samples below ensure the tags are not vacuous.
// =============================================================================

func TestUS_VAL_001_RequestValidation(t *testing.T) {
	// AC1 — Register role ∈ {vendor, customer} and password ≥ 8 chars.
	// (legacy: TestRegisterRequest_Validation)
	t.Run("AC1/register_rejects_short_password_and_bad_role", func(t *testing.T) {
		// --- Given / When / Then ---
		assert.False(t,
			bindsOK(t, `{"email":"x@y.com","password":"short","role":"customer","firstName":"A","lastName":"B"}`,
				&RegisterRequest{}),
			"short password must be rejected")
		assert.False(t,
			bindsOK(t, `{"email":"x@y.com","password":"longenough","role":"superadmin","firstName":"A","lastName":"B"}`,
				&RegisterRequest{}),
			"role must be in {vendor, customer}")
	})

	// AC2 — Cart quantity ∈ [1, 1000].
	// (legacy: TestAddToCartRequest_Validation)
	t.Run("AC2/cart_quantity_in_range", func(t *testing.T) {
		assert.False(t, bindsOK(t, `{"productId":1,"quantity":0}`, &AddToCartRequest{}),
			"quantity 0 must be rejected")
		assert.False(t, bindsOK(t, `{"productId":1,"quantity":1001}`, &AddToCartRequest{}),
			"quantity > 1000 must be rejected")
	})

	// AC3 — Coupon code length 1..50.
	// (legacy: TestApplyCouponRequest_Validation)
	t.Run("AC3/coupon_code_length", func(t *testing.T) {
		long := `"` + repeat("A", 51) + `"`
		assert.False(t, bindsOK(t, `{"promoCode":""}`, &ApplyCouponRequest{}),
			"empty code must be rejected")
		assert.False(t, bindsOK(t, `{"promoCode":`+long+`}`, &ApplyCouponRequest{}),
			"51-char code must be rejected")
	})

	// AC4 — Checkout payment method whitelist.
	// (legacy: TestCheckoutRequest_Validation)
	t.Run("AC4/checkout_payment_method_whitelist", func(t *testing.T) {
		payload := `{"paymentMethod":"crypto","shippingAddress":{"phone":"0123456789","country":"LK","postalCode":"04000","addressLine1":"X"}}`
		assert.False(t, bindsOK(t, payload, &CheckoutRequest{}),
			"unknown payment method must be rejected")
	})

	// AC5 — Checkout promo codes capped at 5.
	// (legacy: TestCheckoutRequest_Validation)
	t.Run("AC5/checkout_promo_codes_capped", func(t *testing.T) {
		payload, _ := json.Marshal(map[string]interface{}{
			"paymentMethod":   "cod",
			"promoCodes":      []string{"A", "B", "C", "D", "E", "F"},
			"shippingAddress": map[string]string{"phone": "0123456789", "country": "LK", "postalCode": "04000", "addressLine1": "X"},
		})
		assert.False(t, bindsOK(t, string(payload), &CheckoutRequest{}),
			"6 promo codes must be rejected (limit 5)")
	})

	// AC6 — Shipping address shape: phone 10, country 2, postal 5.
	// (legacy: TestShippingAddressRequest_Validation)
	t.Run("AC6/shipping_address_shape", func(t *testing.T) {
		assert.False(t, bindsOK(t,
			`{"phone":"123","country":"LK","postalCode":"04000","addressLine1":"X"}`,
			&ShippingAddressRequest{}), "phone must be 10 chars")
		assert.False(t, bindsOK(t,
			`{"phone":"0123456789","country":"LKA","postalCode":"04000","addressLine1":"X"}`,
			&ShippingAddressRequest{}), "country must be 2 chars")
		assert.False(t, bindsOK(t,
			`{"phone":"0123456789","country":"LK","postalCode":"123","addressLine1":"X"}`,
			&ShippingAddressRequest{}), "postalCode must be 5 chars")
	})

	// AC7 — InitiateReturnRequest reason length [10,1000] and ≤5 evidence URLs.
	// (legacy: TestInitiateReturnRequest_Validation)
	t.Run("AC7/return_request_shape", func(t *testing.T) {
		assert.False(t, bindsOK(t,
			`{"orderItemId":1,"reason":"short"}`, &InitiateReturnRequest{}),
			"reason < 10 chars must be rejected")
		evidence := `["a","b","c","d","e","f"]`
		assert.False(t, bindsOK(t,
			`{"orderItemId":1,"reason":"`+repeat("r", 20)+`","evidenceURLs":`+evidence+`}`,
			&InitiateReturnRequest{}),
			"more than 5 evidence URLs must be rejected")
	})
}

// repeat returns s concatenated n times. Kept local to avoid importing "strings"
// for a single use.
func repeat(s string, n int) string {
	out := make([]byte, 0, len(s)*n)
	for i := 0; i < n; i++ {
		out = append(out, s...)
	}
	return string(out)
}
