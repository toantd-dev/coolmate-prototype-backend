package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/datatypes"
)

type Promotion struct {
	ID                   uint
	VendorID             *uint `gorm:"index"` // Null = platform-wide
	Type                 string `gorm:"type:varchar(50);not null"` // product_discount, coupon, bogo, bundle, order_discount, free_shipping, bank_ipg
	Code                 string `gorm:"index"`
	DiscountType         string `gorm:"type:varchar(50);not null"` // flat, percent
	DiscountValue        float64 `gorm:"type:numeric(15,2);not null"`
	MinOrderValue        *float64 `gorm:"type:numeric(15,2)"`
	ValidFrom            time.Time `gorm:"not null"`
	ValidTo              time.Time `gorm:"not null"`
	UsageLimit           *int
	UsedCount            int `gorm:"default:0"`
	FundingType          string `gorm:"type:varchar(50);not null"` // vendor, platform, shared
	VendorSharePct       *float64 `gorm:"type:numeric(5,2)"` // For shared discounts
	ApplicableCategories pq.Int64Array `gorm:"type:int8[]"`
	IsActive             bool `gorm:"default:true"`
	CreatedAt            time.Time
	UpdatedAt            time.Time

	Vendor *Vendor `gorm:"foreignKey:VendorID"`
}

type ProductPromotion struct {
	ID          uint
	ProductID   uint `gorm:"not null;index"`
	PromotionID uint `gorm:"not null;index"`
	CreatedAt   time.Time

	Product   *Product `gorm:"foreignKey:ProductID"`
	Promotion *Promotion `gorm:"foreignKey:PromotionID"`
}

type WalletTransaction struct {
	ID          uint
	VendorID    uint `gorm:"not null;index"`
	OrderID     *uint `gorm:"index"`
	RefType     string `gorm:"type:varchar(50);not null"` // order, return, settlement, adjustment
	Amount      float64 `gorm:"type:numeric(15,2);not null"`
	Type        string `gorm:"type:varchar(50);not null"` // credit, debit
	BalanceAfter float64 `gorm:"type:numeric(15,2);not null"`
	Description string
	CreatedAt   time.Time

	Vendor *Vendor `gorm:"foreignKey:VendorID"`
	Order  *Order `gorm:"foreignKey:OrderID"`
}

type Settlement struct {
	ID           uint
	VendorID     uint `gorm:"not null;index"`
	PeriodStart  time.Time `gorm:"not null;index"`
	PeriodEnd    time.Time `gorm:"not null"`
	GrossSales   float64 `gorm:"type:numeric(15,2);not null"`
	CommissionDeducted float64 `gorm:"type:numeric(15,2);not null"`
	RefundDeductions float64 `gorm:"type:numeric(15,2);default:0"`
	NetPayable   float64 `gorm:"type:numeric(15,2);not null"`
	Status       string `gorm:"type:varchar(50);not null;default:'pending'"` // pending, processing, paid
	InitiatedAt  time.Time
	PaidAt       *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time

	Vendor *Vendor `gorm:"foreignKey:VendorID"`
}

type AuditLog struct {
	ID         uint
	UserID     *uint `gorm:"index"`
	EntityType string `gorm:"type:varchar(100);not null;index"`
	EntityID   uint `gorm:"not null;index"`
	Action     string `gorm:"type:varchar(50);not null;index"`
	BeforeData datatypes.JSONMap `gorm:"type:jsonb"`
	AfterData  datatypes.JSONMap `gorm:"type:jsonb"`
	IP         string
	CreatedAt  time.Time

	User *User `gorm:"foreignKey:UserID"`
}
