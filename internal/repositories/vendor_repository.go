package repositories

import (
	"github.com/coolmate/ecommerce-backend/internal/models"
	"gorm.io/gorm"
)

type IVendorRepository interface {
	Create(vendor *models.Vendor) error
	GetByID(id uint) (*models.Vendor, error)
	GetByUserID(userID uint) (*models.Vendor, error)
	GetBySlug(slug string) (*models.Vendor, error)
	List(status string, limit int, offset int) ([]models.Vendor, int64, error)
	Update(vendor *models.Vendor) error
	UpdateStatus(vendorID uint, status string) error
	CreateDocument(doc *models.VendorDocument) error
	GetDocuments(vendorID uint) ([]models.VendorDocument, error)
	CreateBankDetails(details *models.VendorBankDetails) error
	UpdateBankDetails(details *models.VendorBankDetails) error
	GetBankDetails(vendorID uint) (*models.VendorBankDetails, error)
	GetWallet(vendorID uint) (*models.VendorWallet, error)
	CreateWallet(wallet *models.VendorWallet) error
	UpdateWalletBalance(vendorID uint, amount float64) error
}

type VendorRepository struct {
	db *gorm.DB
}

func NewVendorRepository(db *gorm.DB) IVendorRepository {
	return &VendorRepository{db: db}
}

func (vr *VendorRepository) Create(vendor *models.Vendor) error {
	return vr.db.Create(vendor).Error
}

func (vr *VendorRepository) GetByID(id uint) (*models.Vendor, error) {
	var vendor models.Vendor
	if err := vr.db.Preload("User").Preload("BankDetails").Preload("Documents").Preload("Wallet").First(&vendor, id).Error; err != nil {
		return nil, err
	}
	return &vendor, nil
}

func (vr *VendorRepository) GetByUserID(userID uint) (*models.Vendor, error) {
	var vendor models.Vendor
	if err := vr.db.Preload("User").Preload("BankDetails").Preload("Wallet").Where("user_id = ?", userID).First(&vendor).Error; err != nil {
		return nil, err
	}
	return &vendor, nil
}

func (vr *VendorRepository) GetBySlug(slug string) (*models.Vendor, error) {
	var vendor models.Vendor
	if err := vr.db.Where("store_slug = ?", slug).First(&vendor).Error; err != nil {
		return nil, err
	}
	return &vendor, nil
}

func (vr *VendorRepository) List(status string, limit int, offset int) ([]models.Vendor, int64, error) {
	var vendors []models.Vendor
	var total int64

	query := vr.db
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Model(&models.Vendor{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("User").Preload("Documents").Limit(limit).Offset(offset).Find(&vendors).Error; err != nil {
		return nil, 0, err
	}

	return vendors, total, nil
}

func (vr *VendorRepository) Update(vendor *models.Vendor) error {
	return vr.db.Save(vendor).Error
}

func (vr *VendorRepository) UpdateStatus(vendorID uint, status string) error {
	return vr.db.Model(&models.Vendor{}).Where("id = ?", vendorID).Update("status", status).Error
}

func (vr *VendorRepository) CreateDocument(doc *models.VendorDocument) error {
	return vr.db.Create(doc).Error
}

func (vr *VendorRepository) GetDocuments(vendorID uint) ([]models.VendorDocument, error) {
	var docs []models.VendorDocument
	if err := vr.db.Where("vendor_id = ?", vendorID).Find(&docs).Error; err != nil {
		return nil, err
	}
	return docs, nil
}

func (vr *VendorRepository) CreateBankDetails(details *models.VendorBankDetails) error {
	return vr.db.Create(details).Error
}

func (vr *VendorRepository) UpdateBankDetails(details *models.VendorBankDetails) error {
	return vr.db.Save(details).Error
}

func (vr *VendorRepository) GetBankDetails(vendorID uint) (*models.VendorBankDetails, error) {
	var details models.VendorBankDetails
	if err := vr.db.Where("vendor_id = ?", vendorID).First(&details).Error; err != nil {
		return nil, err
	}
	return &details, nil
}

func (vr *VendorRepository) GetWallet(vendorID uint) (*models.VendorWallet, error) {
	var wallet models.VendorWallet
	if err := vr.db.Where("vendor_id = ?", vendorID).First(&wallet).Error; err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (vr *VendorRepository) CreateWallet(wallet *models.VendorWallet) error {
	return vr.db.Create(wallet).Error
}

func (vr *VendorRepository) UpdateWalletBalance(vendorID uint, amount float64) error {
	return vr.db.Model(&models.VendorWallet{}).Where("vendor_id = ?", vendorID).Update("balance", gorm.Expr("balance + ?", amount)).Error
}
