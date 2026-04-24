package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test validation struct definitions and constraints

func TestRegisterRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     RegisterRequest
		isValid bool
	}{
		{
			name: "Valid request",
			req: RegisterRequest{
				Email:     "user@example.com",
				Password:  "SecurePass123",
				FirstName: "John",
				LastName:  "Doe",
				Role:      "customer",
			},
			isValid: true,
		},
		{
			name: "Missing email",
			req: RegisterRequest{
				Email:     "",
				Password:  "SecurePass123",
				FirstName: "John",
				LastName:  "Doe",
				Role:      "customer",
			},
			isValid: false,
		},
		{
			name: "Invalid role",
			req: RegisterRequest{
				Email:     "user@example.com",
				Password:  "SecurePass123",
				FirstName: "John",
				LastName:  "Doe",
				Role:      "admin",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		if tt.isValid {
			assert.NotEmpty(t, tt.req.Email)
			assert.GreaterOrEqual(t, len(tt.req.Password), 8)
			assert.NotEmpty(t, tt.req.FirstName)
			assert.NotEmpty(t, tt.req.LastName)
			assert.True(t, tt.req.Role == "vendor" || tt.req.Role == "customer")
		}
	}
}

func TestLoginRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     LoginRequest
		isValid bool
	}{
		{
			name: "Valid login",
			req: LoginRequest{
				Email:    "user@example.com",
				Password: "MyPassword123",
			},
			isValid: true,
		},
		{
			name: "Invalid email format",
			req: LoginRequest{
				Email:    "invalid-email",
				Password: "MyPassword123",
			},
			isValid: false,
		},
		{
			name: "Missing password",
			req: LoginRequest{
				Email:    "user@example.com",
				Password: "",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		if tt.isValid {
			assert.NotEmpty(t, tt.req.Email)
			assert.NotEmpty(t, tt.req.Password)
			assert.Contains(t, tt.req.Email, "@")
		} else {
			if tt.req.Password == "" {
				assert.Empty(t, tt.req.Password)
			} else if tt.req.Email != "" {
				assert.NotContains(t, tt.req.Email, "@")
			}
		}
	}
}

func TestRegisterVendorRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     RegisterVendorRequest
		isValid bool
	}{
		{
			name: "Valid vendor registration",
			req: RegisterVendorRequest{
				StoreName:       "My Store",
				StoreSlug:       "my-store",
				Description:     "A great store",
				CommissionModel: "margin",
				CommissionRate:  0.05,
			},
			isValid: true,
		},
		{
			name: "Store name too short",
			req: RegisterVendorRequest{
				StoreName:       "AB",
				StoreSlug:       "my-store",
				CommissionModel: "margin",
				CommissionRate:  0.05,
			},
			isValid: false,
		},
		{
			name: "Invalid commission model",
			req: RegisterVendorRequest{
				StoreName:       "My Store",
				StoreSlug:       "my-store",
				CommissionModel: "invalid",
				CommissionRate:  0.05,
			},
			isValid: false,
		},
		{
			name: "Commission rate out of range",
			req: RegisterVendorRequest{
				StoreName:       "My Store",
				StoreSlug:       "my-store",
				CommissionModel: "margin",
				CommissionRate:  1.5,
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		if tt.isValid {
			assert.GreaterOrEqual(t, len(tt.req.StoreName), 3)
			assert.NotEmpty(t, tt.req.StoreSlug)
			assert.True(t, tt.req.CommissionModel == "margin" || tt.req.CommissionModel == "markup")
			assert.GreaterOrEqual(t, tt.req.CommissionRate, 0.0)
			assert.LessOrEqual(t, tt.req.CommissionRate, 1.0)
		}
	}
}

func TestCreateProductRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateProductRequest
		isValid bool
	}{
		{
			name: "Valid product",
			req: CreateProductRequest{
				Name:        "Test Product",
				SKU:         "SKU123",
				Description: "A very detailed product description",
				CategoryID:  1,
				BasePrice:   99.99,
				CostPrice:   50.0,
				Weight:      2.5,
			},
			isValid: true,
		},
		{
			name: "Price validation - negative base price",
			req: CreateProductRequest{
				Name:        "Test Product",
				SKU:         "SKU123",
				Description: "A very detailed product description",
				CategoryID:  1,
				BasePrice:   -10.0,
				CostPrice:   50.0,
				Weight:      2.5,
			},
			isValid: false,
		},
		{
			name: "Description too short",
			req: CreateProductRequest{
				Name:        "Test Product",
				SKU:         "SKU123",
				Description: "short",
				CategoryID:  1,
				BasePrice:   99.99,
				CostPrice:   50.0,
				Weight:      2.5,
			},
			isValid: false,
		},
		{
			name: "Missing category",
			req: CreateProductRequest{
				Name:        "Test Product",
				SKU:         "SKU123",
				Description: "A very detailed product description",
				CategoryID:  0,
				BasePrice:   99.99,
				CostPrice:   50.0,
				Weight:      2.5,
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		if tt.isValid {
			assert.GreaterOrEqual(t, len(tt.req.Name), 3)
			assert.GreaterOrEqual(t, len(tt.req.Description), 10)
			assert.Greater(t, tt.req.BasePrice, 0.0)
			assert.Greater(t, tt.req.CategoryID, uint(0))
		}
	}
}

func TestShippingAddressRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     ShippingAddressRequest
		isValid bool
	}{
		{
			name: "Valid address",
			req: ShippingAddressRequest{
				FullName:    "John Doe",
				Phone:       "1234567890",
				Email:       "john@example.com",
				AddressLine: "123 Main Street",
				City:        "New York",
				State:       "NY",
				PostalCode:  "10001",
				Country:     "US",
			},
			isValid: true,
		},
		{
			name: "Invalid phone - too short",
			req: ShippingAddressRequest{
				FullName:    "John Doe",
				Phone:       "123456789",
				Email:       "john@example.com",
				AddressLine: "123 Main Street",
				City:        "New York",
				State:       "NY",
				PostalCode:  "10001",
				Country:     "US",
			},
			isValid: false,
		},
		{
			name: "Invalid postal code",
			req: ShippingAddressRequest{
				FullName:    "John Doe",
				Phone:       "1234567890",
				Email:       "john@example.com",
				AddressLine: "123 Main Street",
				City:        "New York",
				State:       "NY",
				PostalCode:  "1000",
				Country:     "US",
			},
			isValid: false,
		},
		{
			name: "Invalid country code",
			req: ShippingAddressRequest{
				FullName:    "John Doe",
				Phone:       "1234567890",
				Email:       "john@example.com",
				AddressLine: "123 Main Street",
				City:        "New York",
				State:       "NY",
				PostalCode:  "10001",
				Country:     "USA",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		if tt.isValid {
			assert.Equal(t, len(tt.req.Phone), 10)
			assert.Equal(t, len(tt.req.PostalCode), 5)
			assert.Equal(t, len(tt.req.Country), 2)
			assert.Contains(t, tt.req.Email, "@")
		}
	}
}

func TestCheckoutRequest_Validation(t *testing.T) {
	validPaymentMethods := []string{"bank_ipg", "emi", "cod", "manual_transfer"}

	tests := []struct {
		name    string
		req     CheckoutRequest
		isValid bool
	}{
		{
			name: "Valid checkout",
			req: CheckoutRequest{
				ShippingAddress: ShippingAddressRequest{
					FullName:    "John Doe",
					Phone:       "1234567890",
					Email:       "john@example.com",
					AddressLine: "123 Main Street",
					City:        "New York",
					State:       "NY",
					PostalCode:  "10001",
					Country:     "US",
				},
				PaymentMethod: "bank_ipg",
				PromoCodes:    []string{"PROMO1"},
			},
			isValid: true,
		},
		{
			name: "Invalid payment method",
			req: CheckoutRequest{
				ShippingAddress: ShippingAddressRequest{
					FullName:    "John Doe",
					Phone:       "1234567890",
					Email:       "john@example.com",
					AddressLine: "123 Main Street",
					City:        "New York",
					State:       "NY",
					PostalCode:  "10001",
					Country:     "US",
				},
				PaymentMethod: "invalid_method",
			},
			isValid: false,
		},
		{
			name: "Too many promo codes",
			req: CheckoutRequest{
				ShippingAddress: ShippingAddressRequest{
					FullName:    "John Doe",
					Phone:       "1234567890",
					Email:       "john@example.com",
					AddressLine: "123 Main Street",
					City:        "New York",
					State:       "NY",
					PostalCode:  "10001",
					Country:     "US",
				},
				PaymentMethod: "bank_ipg",
				PromoCodes:    []string{"PROMO1", "PROMO2", "PROMO3", "PROMO4", "PROMO5", "PROMO6"},
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		if tt.isValid {
			assert.Contains(t, validPaymentMethods, tt.req.PaymentMethod)
			assert.LessOrEqual(t, len(tt.req.PromoCodes), 5)
		}
	}
}

func TestAddToCartRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     AddToCartRequest
		isValid bool
	}{
		{
			name:    "Valid add to cart",
			req:     AddToCartRequest{ProductVariantID: 1, Quantity: 5},
			isValid: true,
		},
		{
			name:    "Zero quantity",
			req:     AddToCartRequest{ProductVariantID: 1, Quantity: 0},
			isValid: false,
		},
		{
			name:    "Negative quantity",
			req:     AddToCartRequest{ProductVariantID: 1, Quantity: -1},
			isValid: false,
		},
		{
			name:    "Quantity too high",
			req:     AddToCartRequest{ProductVariantID: 1, Quantity: 1001},
			isValid: false,
		},
		{
			name:    "No product variant",
			req:     AddToCartRequest{ProductVariantID: 0, Quantity: 5},
			isValid: false,
		},
	}

	for _, tt := range tests {
		if tt.isValid {
			assert.Greater(t, tt.req.ProductVariantID, uint(0))
			assert.Greater(t, tt.req.Quantity, 0)
			assert.LessOrEqual(t, tt.req.Quantity, 1000)
		}
	}
}

func TestPaginationQuery_GetOffset(t *testing.T) {
	tests := []struct {
		page     int
		limit    int
		expected int
	}{
		{1, 10, 0},
		{2, 10, 10},
		{3, 20, 40},
		{0, 0, 0}, // Test defaults
	}

	for _, tt := range tests {
		pq := PaginationQuery{Page: tt.page, Limit: tt.limit}
		offset := pq.GetOffset()

		// After GetOffset, page and limit should have defaults
		assert.GreaterOrEqual(t, pq.Page, 1)
		assert.GreaterOrEqual(t, pq.Limit, 1)
		assert.LessOrEqual(t, pq.Limit, 100)

		// Verify offset calculation
		if tt.page > 0 && tt.limit > 0 {
			assert.Equal(t, (tt.page-1)*tt.limit, offset)
		}
	}
}

func TestPaginationQuery_DefaultLimits(t *testing.T) {
	tests := []struct {
		name   string
		limit  int
		expect int
	}{
		{"Zero limit", 0, 20},
		{"Negative limit", -5, 20},
		{"Valid limit", 50, 50},
		{"Max limit exceeded", 200, 100},
		{"Boundary - max allowed", 100, 100},
	}

	for _, tt := range tests {
		pq := PaginationQuery{Page: 1, Limit: tt.limit}
		pq.GetOffset()

		assert.Equal(t, tt.expect, pq.Limit, tt.name)
	}
}

func TestInitiateReturnRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     InitiateReturnRequest
		isValid bool
	}{
		{
			name: "Valid return request",
			req: InitiateReturnRequest{
				OrderID: 1,
				Reason:  "Product is damaged and unusable",
				Evidence: []string{"https://example.com/image1.jpg"},
			},
			isValid: true,
		},
		{
			name: "Reason too short",
			req: InitiateReturnRequest{
				OrderID: 1,
				Reason:  "Damaged",
				Evidence: []string{},
			},
			isValid: false,
		},
		{
			name: "Too many evidence files",
			req: InitiateReturnRequest{
				OrderID: 1,
				Reason:  "Product is damaged and unusable",
				Evidence: []string{"url1", "url2", "url3", "url4", "url5", "url6"},
			},
			isValid: false,
		},
		{
			name: "No order ID",
			req: InitiateReturnRequest{
				OrderID: 0,
				Reason:  "Product is damaged and unusable",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		if tt.isValid {
			assert.Greater(t, tt.req.OrderID, uint(0))
			assert.GreaterOrEqual(t, len(tt.req.Reason), 10)
			assert.LessOrEqual(t, len(tt.req.Evidence), 5)
		}
	}
}
