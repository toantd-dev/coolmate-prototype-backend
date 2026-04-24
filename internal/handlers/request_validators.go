package handlers

// Auth Requests
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required,min=2"`
	LastName  string `json:"last_name" binding:"required,min=2"`
	Role      string `json:"role" binding:"required,oneof=vendor customer"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Vendor Requests
type RegisterVendorRequest struct {
	StoreName        string `json:"store_name" binding:"required,min=3,max=100"`
	StoreSlug        string `json:"store_slug" binding:"required,min=3,max=50,alphanum"`
	Description      string `json:"description" binding:"max=1000"`
	CommissionModel  string `json:"commission_model" binding:"required,oneof=margin markup"`
	CommissionRate   float64 `json:"commission_rate" binding:"required,min=0,max=1"`
}

type UpdateVendorProfileRequest struct {
	StoreName   string `json:"store_name" binding:"min=3,max=100"`
	Description string `json:"description" binding:"max=1000"`
	LogoURL     string `json:"logo_url" binding:"url"`
}

type UpdateBankDetailsRequest struct {
	AccountName   string `json:"account_name" binding:"required,min=5"`
	AccountNumber string `json:"account_number" binding:"required,min=8"`
	BankName      string `json:"bank_name" binding:"required,min=2"`
	Branch        string `json:"branch" binding:"required,min=2"`
}

// Product Requests
type CreateProductRequest struct {
	Name             string  `json:"name" binding:"required,min=3,max=255"`
	SKU              string  `json:"sku" binding:"required,min=3,max=50"`
	Description      string  `json:"description" binding:"required,min=10,max=5000"`
	CategoryID       uint    `json:"category_id" binding:"required,gt=0"`
	BrandID          uint    `json:"brand_id" binding:"gt=0"`
	BasePrice        float64 `json:"base_price" binding:"required,gt=0"`
	CostPrice        float64 `json:"cost_price" binding:"required,gte=0"`
	Weight           float64 `json:"weight" binding:"required,gt=0"`
	IsReturnable     bool    `json:"is_returnable"`
	ReturnWindowDays int     `json:"return_window_days"`
	Warranty         string  `json:"warranty" binding:"max=500"`
}

type UpdateProductRequest struct {
	Name             string  `json:"name" binding:"min=3,max=255"`
	Description      string  `json:"description" binding:"min=10,max=5000"`
	BasePrice        float64 `json:"base_price" binding:"gt=0"`
	CostPrice        float64 `json:"cost_price" binding:"gte=0"`
	Weight           float64 `json:"weight" binding:"gt=0"`
	IsReturnable     bool    `json:"is_returnable"`
	ReturnWindowDays int     `json:"return_window_days"`
}

type CreateProductVariantRequest struct {
	SKU       string                 `json:"sku" binding:"required,min=3,max=50"`
	Price     float64                `json:"price" binding:"required,gt=0"`
	Stock     int                    `json:"stock" binding:"required,gte=0"`
	Attributes map[string]interface{} `json:"attributes" binding:"required"`
}

// Order Requests
type AddToCartRequest struct {
	ProductVariantID uint `json:"product_variant_id" binding:"required,gt=0"`
	Quantity         int  `json:"quantity" binding:"required,gt=0,max=1000"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,gt=0,max=1000"`
}

type CheckoutRequest struct {
	ShippingAddress ShippingAddressRequest `json:"shipping_address" binding:"required"`
	PaymentMethod   string                 `json:"payment_method" binding:"required,oneof=bank_ipg emi cod manual_transfer"`
	PromoCodes      []string               `json:"promo_codes" binding:"max=5"`
}

type ShippingAddressRequest struct {
	FullName    string `json:"full_name" binding:"required,min=3"`
	Phone       string `json:"phone" binding:"required,len=10"`
	Email       string `json:"email" binding:"required,email"`
	AddressLine string `json:"address_line" binding:"required,min=5"`
	City        string `json:"city" binding:"required,min=2"`
	State       string `json:"state" binding:"required,min=2"`
	PostalCode  string `json:"postal_code" binding:"required,len=5"`
	Country     string `json:"country" binding:"required,len=2"`
}

type ApplyCouponRequest struct {
	PromoCode string `json:"promo_code" binding:"required,min=3,max=50"`
}

// Return Requests
type InitiateReturnRequest struct {
	OrderID uint     `json:"order_id" binding:"required,gt=0"`
	Reason  string   `json:"reason" binding:"required,min=10,max=1000"`
	Evidence []string `json:"evidence_urls" binding:"max=5"`
}

// Promotion Requests
type CreatePromotionRequest struct {
	Type             string   `json:"type" binding:"required,oneof=product_discount coupon bogo bundle order_discount free_shipping bank_ipg"`
	Code             string   `json:"code" binding:"max=50"`
	DiscountType     string   `json:"discount_type" binding:"required,oneof=flat percent"`
	DiscountValue    float64  `json:"discount_value" binding:"required,gt=0"`
	MinOrderValue    float64  `json:"min_order_value" binding:"gte=0"`
	ValidFrom        string   `json:"valid_from" binding:"required,datetime=2006-01-02T15:04:05Z07:00"`
	ValidTo          string   `json:"valid_to" binding:"required,datetime=2006-01-02T15:04:05Z07:00"`
	UsageLimit       int      `json:"usage_limit" binding:"gte=0"`
	ApplicableCategories []uint `json:"applicable_categories"`
}

// Pagination
type PaginationQuery struct {
	Page  int `form:"page" binding:"min=1" default:"1"`
	Limit int `form:"limit" binding:"min=1,max=100" default:"20"`
}

func (p *PaginationQuery) GetOffset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = 20
	} else if p.Limit > 100 {
		p.Limit = 100
	}
	return (p.Page - 1) * p.Limit
}
