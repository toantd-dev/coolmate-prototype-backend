package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/coolmate/ecommerce-backend/internal/models"
	"github.com/coolmate/ecommerce-backend/internal/repositories"
	"github.com/coolmate/ecommerce-backend/pkg/cache"
	"github.com/coolmate/ecommerce-backend/pkg/storage"
)

type VendorService struct {
	vendorRepo repositories.IVendorRepository
	userRepo   repositories.IUserRepository
	s3Manager  *storage.S3Manager
	cache      *cache.CacheManager
}

func NewVendorService(
	vendorRepo repositories.IVendorRepository,
	userRepo repositories.IUserRepository,
	s3Manager *storage.S3Manager,
	cache *cache.CacheManager,
) *VendorService {
	return &VendorService{
		vendorRepo: vendorRepo,
		userRepo:   userRepo,
		s3Manager:  s3Manager,
		cache:      cache,
	}
}

func (vs *VendorService) GetVendorByID(id uint) (*models.Vendor, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("vendor:%d", id)

	// Try cache first
	vendor := &models.Vendor{}
	err := vs.cache.Get(ctx, cacheKey, vendor)
	if err == nil {
		return vendor, nil
	}

	vendor, err = vs.vendorRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Cache for 1 hour
	vs.cache.Set(ctx, cacheKey, vendor, time.Hour)
	return vendor, nil
}

func (vs *VendorService) ListVendors(status string, limit int, offset int) ([]models.Vendor, int64, error) {
	return vs.vendorRepo.List(status, limit, offset)
}

func (vs *VendorService) ApproveVendor(vendorID uint) error {
	ctx := context.Background()
	if err := vs.vendorRepo.UpdateStatus(vendorID, "approved"); err != nil {
		return err
	}
	vs.cache.Delete(ctx, fmt.Sprintf("vendor:%d", vendorID))
	return nil
}

func (vs *VendorService) RejectVendor(vendorID uint) error {
	ctx := context.Background()
	if err := vs.vendorRepo.UpdateStatus(vendorID, "rejected"); err != nil {
		return err
	}
	vs.cache.Delete(ctx, fmt.Sprintf("vendor:%d", vendorID))
	return nil
}

func (vs *VendorService) SuspendVendor(vendorID uint) error {
	ctx := context.Background()
	if err := vs.vendorRepo.UpdateStatus(vendorID, "suspended"); err != nil {
		return err
	}
	vs.cache.Delete(ctx, fmt.Sprintf("vendor:%d", vendorID))
	return nil
}

// CanListProducts checks if vendor is eligible to list products
// Requirements: approved status + agreed to latest agreement
func (vs *VendorService) CanListProducts(vendorID uint) (bool, error) {
	vendor, err := vs.GetVendorByID(vendorID)
	if err != nil {
		return false, err
	}

	if vendor.Status != "approved" {
		return false, fmt.Errorf("vendor status is %s, not approved", vendor.Status)
	}

	// Check if vendor has accepted latest agreement
	if vendor.AgreementAcceptedAt == nil {
		return false, errors.New("vendor has not agreed to terms")
	}

	return true, nil
}

// ValidateVendor checks vendor data integrity
func (vs *VendorService) ValidateVendor(vendor *models.Vendor) error {
	if vendor == nil {
		return errors.New("vendor cannot be nil")
	}

	if vendor.UserID == 0 {
		return errors.New("user ID is required")
	}

	if vendor.StoreName == "" || len(vendor.StoreName) < 3 {
		return errors.New("store name must be at least 3 characters")
	}

	if vendor.StoreSlug == "" {
		return errors.New("store slug is required")
	}

	if vendor.CommissionModel != "margin" && vendor.CommissionModel != "markup" {
		return errors.New("commission model must be 'margin' or 'markup'")
	}

	if vendor.CommissionRate < 0 || vendor.CommissionRate > 1 {
		return errors.New("commission rate must be between 0 and 1")
	}

	return nil
}

// GetVendorWallet retrieves vendor's wallet with transaction history
func (vs *VendorService) GetVendorWallet(vendorID uint) (*models.VendorWallet, error) {
	return vs.vendorRepo.GetWallet(vendorID)
}

// UpdateBankDetails updates vendor bank details (admin only)
func (vs *VendorService) UpdateBankDetails(
	vendorID uint,
	accountName string,
	accountNumber string,
	bankName string,
	branchName string,
) error {
	ctx := context.Background()
	details := &models.VendorBankDetails{
		VendorID:      vendorID,
		AccountName:   accountName,
		AccountNumber: accountNumber,
		BankName:      bankName,
		BranchName:    branchName,
	}

	if err := vs.vendorRepo.UpdateBankDetails(details); err != nil {
		return err
	}

	vs.cache.Delete(ctx, fmt.Sprintf("vendor:%d", vendorID))
	return nil
}
