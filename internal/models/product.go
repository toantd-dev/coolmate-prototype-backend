package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Category struct {
	ID               uint
	Name             string `gorm:"not null;index"`
	Slug             string `gorm:"uniqueIndex;not null"`
	ParentID         *uint `gorm:"index"`
	MinPrice         *float64 `gorm:"type:numeric(15,2)"`
	MaxPrice         *float64 `gorm:"type:numeric(15,2)"`
	MaxDiscountPct   *float64 `gorm:"type:numeric(5,2)"`
	CommissionRate   *float64 `gorm:"type:numeric(5,2)"`
	CommissionModel  string `gorm:"type:varchar(50)"` // markup, margin
	IsActive         bool `gorm:"default:true"`
	CreatedAt        time.Time
	UpdatedAt        time.Time

	SubCategories []Category `gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Products      []Product `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Brand struct {
	ID        uint
	Name      string `gorm:"not null;uniqueIndex"`
	Slug      string `gorm:"uniqueIndex;not null"`
	LogoURL   string
	IsActive  bool `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Products []Product `gorm:"foreignKey:BrandID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Product struct {
	ID             uint
	VendorID       uint `gorm:"not null;index"`
	CategoryID     uint `gorm:"not null;index"`
	BrandID        *uint `gorm:"index"`
	Name           string `gorm:"not null;index"`
	Slug           string `gorm:"uniqueIndex;not null"`
	SKU            string `gorm:"uniqueIndex;not null"`
	Description    string `gorm:"type:text"`
	BasePrice      float64 `gorm:"type:numeric(15,2);not null"`
	CostPrice      *float64 `gorm:"type:numeric(15,2)"` // Only visible to admin
	Status         string `gorm:"type:varchar(50);not null;default:'draft'"` // draft, pending_approval, published, archived, rejected
	Weight         *float64
	Dimensions     datatypes.JSONMap `gorm:"type:jsonb"`
	Warranty       *string
	SEOTitle       string
	SEODescription string
	IsReturnable   bool `gorm:"default:true"`
	ReturnWindowDays int `gorm:"default:30"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`

	// Relations
	Vendor         *Vendor `gorm:"foreignKey:VendorID"`
	Category       *Category `gorm:"foreignKey:CategoryID"`
	Brand          *Brand `gorm:"foreignKey:BrandID"`
	Variants       []ProductVariant `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Images         []ProductImage `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ApprovalLogs   []ProductApprovalLog `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type ProductVariant struct {
	ID         uint
	ProductID  uint `gorm:"not null;index"`
	SKU        string `gorm:"not null;uniqueIndex:,composite:product_id"`
	Price      float64 `gorm:"type:numeric(15,2);not null"`
	Stock      int `gorm:"not null;default:0"`
	Attributes datatypes.JSONMap `gorm:"type:jsonb"` // {color: red, size: M}
	IsActive   bool `gorm:"default:true"`
	CreatedAt  time.Time
	UpdatedAt  time.Time

	Product *Product `gorm:"foreignKey:ProductID"`
}

type ProductImage struct {
	ID        uint
	ProductID uint `gorm:"not null;index"`
	VariantID *uint `gorm:"index"`
	URL       string `gorm:"not null"`
	IsPrimary bool `gorm:"default:false"`
	SortOrder int
	CreatedAt time.Time
	UpdatedAt time.Time

	Product *Product `gorm:"foreignKey:ProductID"`
	Variant *ProductVariant `gorm:"foreignKey:VariantID"`
}

type ProductApprovalLog struct {
	ID        uint
	ProductID uint `gorm:"not null;index"`
	AdminID   uint `gorm:"not null;index"`
	Action    string `gorm:"type:varchar(50);not null"` // approved, rejected
	Comment   string
	CreatedAt time.Time

	Product *Product `gorm:"foreignKey:ProductID"`
	Admin   *User `gorm:"foreignKey:AdminID"`
}
