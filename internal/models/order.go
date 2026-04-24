package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Cart struct {
	ID        uint
	UserID    *uint `gorm:"index"`
	SessionID string `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time

	User  *User `gorm:"foreignKey:UserID"`
	Items []CartItem `gorm:"foreignKey:CartID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type CartItem struct {
	ID        uint
	CartID    uint `gorm:"not null;index"`
	ProductID uint `gorm:"not null;index"`
	VariantID *uint `gorm:"index"`
	Quantity  int `gorm:"not null;default:1"`
	UnitPrice float64 `gorm:"type:numeric(15,2);not null"`
	AddedAt   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time

	Cart      *Cart `gorm:"foreignKey:CartID"`
	Product   *Product `gorm:"foreignKey:ProductID"`
	Variant   *ProductVariant `gorm:"foreignKey:VariantID"`
}

type Order struct {
	ID                 uint
	CustomerID         *uint `gorm:"index"`
	GuestEmail         string
	Status             string `gorm:"type:varchar(50);not null;default:'pending'"` // pending, paid, processing, shipped, delivered, cancelled, refunded
	PaymentMethod      string `gorm:"type:varchar(50);not null"` // bank_ipg, emi, cod, manual_transfer
	PaymentStatus      string `gorm:"type:varchar(50);not null;default:'pending'"` // pending, paid, failed, refunded
	Subtotal           float64 `gorm:"type:numeric(15,2);not null"`
	DiscountTotal      float64 `gorm:"type:numeric(15,2);default:0"`
	ShippingTotal      float64 `gorm:"type:numeric(15,2);default:0"`
	GrandTotal         float64 `gorm:"type:numeric(15,2);not null"`
	ShippingAddress    datatypes.JSONMap `gorm:"type:jsonb"`
	PromoCodesApplied  pq.StringArray `gorm:"type:text[]"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          gorm.DeletedAt `gorm:"index"`

	Customer  *User `gorm:"foreignKey:CustomerID"`
	Items     []OrderItem `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	SubOrders []SubOrder `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	StatusLogs []OrderStatusLog `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type SubOrder struct {
	ID                 uint
	OrderID            uint `gorm:"not null;index"`
	VendorID           uint `gorm:"not null;index"`
	Status             string `gorm:"type:varchar(50);not null;default:'pending'"` // pending, ready_to_ship, shipped, delivered, cancelled
	Subtotal           float64 `gorm:"type:numeric(15,2);not null"`
	CommissionAmount   float64 `gorm:"type:numeric(15,2);default:0"`
	VendorEarning      float64 `gorm:"type:numeric(15,2);not null"`
	ShippingCharge     float64 `gorm:"type:numeric(15,2);default:0"`
	CreatedAt          time.Time
	UpdatedAt          time.Time

	Order  *Order `gorm:"foreignKey:OrderID"`
	Vendor *Vendor `gorm:"foreignKey:VendorID"`
}

type OrderItem struct {
	ID                    uint
	OrderID               uint `gorm:"not null;index"`
	SubOrderID            uint `gorm:"not null;index"`
	ProductID             uint `gorm:"not null;index"`
	VariantID             *uint `gorm:"index"`
	VendorID              uint `gorm:"not null;index"`
	Quantity              int `gorm:"not null"`
	UnitPrice             float64 `gorm:"type:numeric(15,2);not null"`
	DiscountAmount        float64 `gorm:"type:numeric(15,2);default:0"`
	CommissionAmount      float64 `gorm:"type:numeric(15,2);default:0"`
	VendorEarning         float64 `gorm:"type:numeric(15,2);not null"`
	CommissionModelSnapshot datatypes.JSONMap `gorm:"type:jsonb"`
	CreatedAt             time.Time
	UpdatedAt             time.Time

	Order   *Order `gorm:"foreignKey:OrderID"`
	SubOrder *SubOrder `gorm:"foreignKey:SubOrderID"`
	Product *Product `gorm:"foreignKey:ProductID"`
	Variant *ProductVariant `gorm:"foreignKey:VariantID"`
	Vendor  *Vendor `gorm:"foreignKey:VendorID"`
}

type OrderStatusLog struct {
	ID         uint
	OrderID    uint `gorm:"not null;index"`
	SubOrderID *uint `gorm:"index"`
	ChangedBy  uint `gorm:"not null;index"`
	FromStatus string
	ToStatus   string
	Note       string
	CreatedAt  time.Time

	Order   *Order `gorm:"foreignKey:OrderID"`
	SubOrder *SubOrder `gorm:"foreignKey:SubOrderID"`
	User    *User `gorm:"foreignKey:ChangedBy"`
}

type ReturnRequest struct {
	ID           uint
	OrderID      uint `gorm:"not null;index"`
	SubOrderID   uint `gorm:"not null;index"`
	CustomerID   uint `gorm:"not null;index"`
	Reason       string `gorm:"type:text;not null"`
	EvidenceURLs pq.StringArray `gorm:"type:text[]"`
	Status       string `gorm:"type:varchar(50);not null;default:'requested'"` // requested, approved, rejected, picked_up, completed
	VendorNote   string
	AdminNote    string
	CreatedAt    time.Time
	UpdatedAt    time.Time

	Order     *Order `gorm:"foreignKey:OrderID"`
	SubOrder  *SubOrder `gorm:"foreignKey:SubOrderID"`
	Customer  *User `gorm:"foreignKey:CustomerID"`
	Refund    *Refund `gorm:"foreignKey:ReturnRequestID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Refund struct {
	ID               uint
	ReturnRequestID  uint `gorm:"not null;index"`
	OrderID          uint `gorm:"not null;index"`
	Amount           float64 `gorm:"type:numeric(15,2);not null"`
	Method           string `gorm:"type:varchar(50);not null"` // original, bank_transfer
	Status           string `gorm:"type:varchar(50);not null;default:'pending'"` // pending, processed
	ProcessedAt      *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time

	ReturnRequest *ReturnRequest `gorm:"foreignKey:ReturnRequestID"`
	Order         *Order `gorm:"foreignKey:OrderID"`
}
