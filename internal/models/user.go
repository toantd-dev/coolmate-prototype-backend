package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleSuperAdmin UserRole = "super_admin"
	RoleAdmin      UserRole = "admin"
	RoleVendor     UserRole = "vendor"
	RoleCustomer   UserRole = "customer"
)

type UserStatus string

const (
	UserActive   UserStatus = "active"
	UserInactive UserStatus = "inactive"
	UserBanned   UserStatus = "banned"
)

type User struct {
	ID               uint
	Email            string `gorm:"uniqueIndex;not null"`
	PasswordHash     string `gorm:"not null"`
	FirstName        string
	LastName         string
	Phone            string
	Role             UserRole `gorm:"type:varchar(50);not null;default:'customer'"`
	Status           UserStatus `gorm:"type:varchar(50);not null;default:'active'"`
	EmailVerifiedAt  *time.Time
	EmailVerifyToken string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`

	// Relations
	Vendor        *Vendor `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	RefreshTokens []RefreshToken `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type RefreshToken struct {
	ID        uint
	UserID    uint `gorm:"not null;index"`
	TokenHash string `gorm:"not null;index"`
	ExpiresAt time.Time `gorm:"not null"`
	RevokedAt *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time

	User *User `gorm:"foreignKey:UserID"`
}

type Vendor struct {
	ID                 uint
	UserID             uint `gorm:"uniqueIndex;not null"`
	VendorType         string `gorm:"type:varchar(50);not null"` // individual, business
	StoreName          string `gorm:"not null"`
	StoreSlug          string `gorm:"uniqueIndex;not null"`
	LogoURL            string
	Description        string
	Status             string `gorm:"type:varchar(50);not null;default:'pending'"` // pending, approved, suspended, rejected
	CommissionModel    string `gorm:"type:varchar(50);not null;default:'markup'"` // markup, margin
	CommissionRate     float64 `gorm:"type:numeric(5,2)"`
	AgreementAcceptedAt *time.Time
	AgreementVersion   int `gorm:"default:0"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          gorm.DeletedAt `gorm:"index"`

	// Relations
	User              *User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	BankDetails       *VendorBankDetails `gorm:"foreignKey:VendorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Documents         []VendorDocument `gorm:"foreignKey:VendorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Staff             []VendorStaff `gorm:"foreignKey:VendorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Products          []Product `gorm:"foreignKey:VendorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Wallet            *VendorWallet `gorm:"foreignKey:VendorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type VendorBankDetails struct {
	ID            uint
	VendorID      uint `gorm:"uniqueIndex;not null"`
	AccountName   string
	AccountNumber string
	BankName      string
	BranchName    string
	CreatedAt     time.Time
	UpdatedAt     time.Time

	Vendor *Vendor `gorm:"foreignKey:VendorID"`
}

type VendorDocument struct {
	ID         uint
	VendorID   uint `gorm:"not null;index"`
	DocType    string `gorm:"type:varchar(100);not null"` // nic, passport, br_cert, tin, bank_proof
	FileURL    string `gorm:"not null"`
	Status     string `gorm:"type:varchar(50);not null;default:'pending'"` // pending, verified, rejected
	UploadedAt time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time

	Vendor *Vendor `gorm:"foreignKey:VendorID"`
}

type VendorStaff struct {
	ID        uint
	VendorID  uint `gorm:"not null;index"`
	UserID    uint `gorm:"not null;index"`
	Role      string `gorm:"type:varchar(50);not null"` // store_manager, staff
	IsActive  bool `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Vendor *Vendor `gorm:"foreignKey:VendorID"`
	User   *User `gorm:"foreignKey:UserID"`
}

type VendorWallet struct {
	ID             uint
	VendorID       uint `gorm:"uniqueIndex;not null"`
	Balance        float64 `gorm:"type:numeric(15,2);default:0"`
	PendingBalance float64 `gorm:"type:numeric(15,2);default:0"`
	CreatedAt      time.Time
	UpdatedAt      time.Time

	Vendor *Vendor `gorm:"foreignKey:VendorID"`
}

type VendorAgreement struct {
	ID        uint
	Version   int `gorm:"not null;uniqueIndex"`
	Title     string `gorm:"not null"`
	FileURL   string `gorm:"not null"`
	IsActive  bool `gorm:"not null;default:true"`
	PublishedAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time

	Acceptances []VendorAgreementAcceptance `gorm:"foreignKey:AgreementID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type VendorAgreementAcceptance struct {
	ID          uint
	VendorID    uint `gorm:"not null;index"`
	AgreementID uint `gorm:"not null;index"`
	AcceptedAt  time.Time `gorm:"not null"`
	IPAddress   string
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Vendor    *Vendor `gorm:"foreignKey:VendorID"`
	Agreement *VendorAgreement `gorm:"foreignKey:AgreementID"`
}
