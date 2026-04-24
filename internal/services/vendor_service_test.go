package services

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/coolmate/ecommerce-backend/internal/models"
	"github.com/coolmate/ecommerce-backend/pkg/cache"
	"github.com/coolmate/ecommerce-backend/pkg/storage"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) UpdatePassword(userID uint, passwordHash string) error {
	args := m.Called(userID, passwordHash)
	return args.Error(0)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetRefreshToken(tokenHash string) (*models.RefreshToken, error) {
	args := m.Called(tokenHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RefreshToken), args.Error(1)
}

func (m *MockUserRepository) SaveRefreshToken(token *models.RefreshToken) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockUserRepository) RevokeRefreshToken(tokenHash string) error {
	args := m.Called(tokenHash)
	return args.Error(0)
}

// Test GetVendorByID

func TestGetVendorByID_FromCache(t *testing.T) {
	vendorRepo := new(MockVendorRepository)
	userRepo := new(MockUserRepository)
	cacheManager := cache.NewCacheManager(nil)
	s3Manager := &storage.S3Manager{}

	service := NewVendorService(vendorRepo, userRepo, s3Manager, cacheManager)

	vendor := &models.Vendor{ID: 1, StoreName: "Test Store"}
	vendorRepo.On("GetByID", uint(1)).Return(vendor, nil)

	result, err := service.GetVendorByID(1)

	require.NoError(t, err)
	assert.Equal(t, vendor, result)
	vendorRepo.AssertExpectations(t)
}

func TestGetVendorByID_FromDatabase(t *testing.T) {
	vendorRepo := new(MockVendorRepository)
	userRepo := new(MockUserRepository)
	cacheManager := cache.NewCacheManager(nil)
	s3Manager := &storage.S3Manager{}

	service := NewVendorService(vendorRepo, userRepo, s3Manager, cacheManager)

	vendor := &models.Vendor{ID: 1, StoreName: "Test Store"}
	vendorRepo.On("GetByID", uint(1)).Return(vendor, nil)

	result, err := service.GetVendorByID(1)

	require.NoError(t, err)
	assert.Equal(t, vendor, result)
	vendorRepo.AssertExpectations(t)
}

func TestGetVendorByID_NotFound(t *testing.T) {
	vendorRepo := new(MockVendorRepository)
	userRepo := new(MockUserRepository)
	cacheManager := cache.NewCacheManager(nil)
	s3Manager := &storage.S3Manager{}

	service := NewVendorService(vendorRepo, userRepo, s3Manager, cacheManager)

	vendorRepo.On("GetByID", uint(999)).Return(nil, errors.New("not found"))

	result, err := service.GetVendorByID(999)

	require.Error(t, err)
	assert.Nil(t, result)
}

// Test ListVendors

func TestListVendors_Success(t *testing.T) {
	vendorRepo := new(MockVendorRepository)
	userRepo := new(MockUserRepository)
	cacheManager := cache.NewCacheManager(nil)
	s3Manager := &storage.S3Manager{}

	service := NewVendorService(vendorRepo, userRepo, s3Manager, cacheManager)

	vendors := []models.Vendor{
		{ID: 1, StoreName: "Store 1"},
	}

	vendorRepo.On("List", "approved", 10, 0).Return(vendors, int64(1), nil)

	result, total, err := service.ListVendors("approved", 10, 0)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), total)
	vendorRepo.AssertExpectations(t)
}

// Test ApproveVendor

func TestApproveVendor_Success(t *testing.T) {
	vendorRepo := new(MockVendorRepository)
	userRepo := new(MockUserRepository)
	cacheManager := cache.NewCacheManager(nil)
	s3Manager := &storage.S3Manager{}

	service := NewVendorService(vendorRepo, userRepo, s3Manager, cacheManager)

	vendorRepo.On("UpdateStatus", uint(1), "approved").Return(nil)

	err := service.ApproveVendor(1)

	require.NoError(t, err)
	vendorRepo.AssertExpectations(t)
}

func TestApproveVendor_UpdateFailed(t *testing.T) {
	vendorRepo := new(MockVendorRepository)
	userRepo := new(MockUserRepository)
	cacheManager := cache.NewCacheManager(nil)
	s3Manager := &storage.S3Manager{}

	service := NewVendorService(vendorRepo, userRepo, s3Manager, cacheManager)

	vendorRepo.On("UpdateStatus", uint(999), "approved").Return(errors.New("not found"))

	err := service.ApproveVendor(999)

	require.Error(t, err)
}

// Test RejectVendor

func TestRejectVendor_Success(t *testing.T) {
	vendorRepo := new(MockVendorRepository)
	userRepo := new(MockUserRepository)
	cacheManager := cache.NewCacheManager(nil)
	s3Manager := &storage.S3Manager{}

	service := NewVendorService(vendorRepo, userRepo, s3Manager, cacheManager)

	vendorRepo.On("UpdateStatus", uint(1), "rejected").Return(nil)

	err := service.RejectVendor(1)

	require.NoError(t, err)
	vendorRepo.AssertExpectations(t)
}

// Test SuspendVendor

func TestSuspendVendor_Success(t *testing.T) {
	vendorRepo := new(MockVendorRepository)
	userRepo := new(MockUserRepository)
	cacheManager := cache.NewCacheManager(nil)
	s3Manager := &storage.S3Manager{}

	service := NewVendorService(vendorRepo, userRepo, s3Manager, cacheManager)

	vendorRepo.On("UpdateStatus", uint(1), "suspended").Return(nil)

	err := service.SuspendVendor(1)

	require.NoError(t, err)
	vendorRepo.AssertExpectations(t)
}

// Test CanListProducts

func TestCanListProducts_Approved(t *testing.T) {
	vendorRepo := new(MockVendorRepository)
	userRepo := new(MockUserRepository)
	cacheManager := cache.NewCacheManager(nil)
	s3Manager := &storage.S3Manager{}

	service := NewVendorService(vendorRepo, userRepo, s3Manager, cacheManager)

	agreementTime := time.Now()
	vendor := &models.Vendor{
		ID:       1,
		Status:   "approved",
		AgreementAcceptedAt: &agreementTime,
	}

	vendorRepo.On("GetByID", uint(1)).Return(vendor, nil)

	canList, err := service.CanListProducts(1)

	require.NoError(t, err)
	assert.True(t, canList)
}

func TestCanListProducts_NotApproved(t *testing.T) {
	vendorRepo := new(MockVendorRepository)
	userRepo := new(MockUserRepository)
	cacheManager := cache.NewCacheManager(nil)
	s3Manager := &storage.S3Manager{}

	service := NewVendorService(vendorRepo, userRepo, s3Manager, cacheManager)

	vendor := &models.Vendor{
		ID:     1,
		Status: "pending",
	}

	vendorRepo.On("GetByID", uint(1)).Return(vendor, nil)

	canList, err := service.CanListProducts(1)

	require.Error(t, err)
	assert.False(t, canList)
}

func TestCanListProducts_NoAgreement(t *testing.T) {
	vendorRepo := new(MockVendorRepository)
	userRepo := new(MockUserRepository)
	cacheManager := cache.NewCacheManager(nil)
	s3Manager := &storage.S3Manager{}

	service := NewVendorService(vendorRepo, userRepo, s3Manager, cacheManager)

	vendor := &models.Vendor{
		ID:       1,
		Status:   "approved",
		AgreementAcceptedAt: nil,
	}

	vendorRepo.On("GetByID", uint(1)).Return(vendor, nil)

	canList, err := service.CanListProducts(1)

	require.Error(t, err)
	assert.False(t, canList)
}

func TestCanListProducts_VendorNotFound(t *testing.T) {
	vendorRepo := new(MockVendorRepository)
	userRepo := new(MockUserRepository)
	cacheManager := cache.NewCacheManager(nil)
	s3Manager := &storage.S3Manager{}

	service := NewVendorService(vendorRepo, userRepo, s3Manager, cacheManager)

	vendorRepo.On("GetByID", uint(999)).Return(nil, errors.New("not found"))

	canList, err := service.CanListProducts(999)

	require.Error(t, err)
	assert.False(t, canList)
}

// Test ValidateVendor

func TestValidateVendor_Valid(t *testing.T) {
	vendorRepo := new(MockVendorRepository)
	userRepo := new(MockUserRepository)
	cacheManager := cache.NewCacheManager(nil)
	s3Manager := &storage.S3Manager{}

	service := NewVendorService(vendorRepo, userRepo, s3Manager, cacheManager)

	vendor := &models.Vendor{
		UserID:          1,
		StoreName:       "Test Store",
		StoreSlug:       "test-store",
		CommissionModel: "margin",
		CommissionRate:  0.05,
	}

	err := service.ValidateVendor(vendor)
	assert.NoError(t, err)
}

func TestValidateVendor_NilVendor(t *testing.T) {
	vendorRepo := new(MockVendorRepository)
	userRepo := new(MockUserRepository)
	cacheManager := cache.NewCacheManager(nil)
	s3Manager := &storage.S3Manager{}

	service := NewVendorService(vendorRepo, userRepo, s3Manager, cacheManager)

	err := service.ValidateVendor(nil)
	assert.Error(t, err)
}

func TestValidateVendor_MissingUserID(t *testing.T) {
	vendorRepo := new(MockVendorRepository)
	userRepo := new(MockUserRepository)
	cacheManager := cache.NewCacheManager(nil)
	s3Manager := &storage.S3Manager{}

	service := NewVendorService(vendorRepo, userRepo, s3Manager, cacheManager)

	vendor := &models.Vendor{
		UserID:    0,
		StoreName: "Test Store",
	}

	err := service.ValidateVendor(vendor)
	assert.Error(t, err)
}

func TestValidateVendor_InvalidStoreName(t *testing.T) {
	vendorRepo := new(MockVendorRepository)
	userRepo := new(MockUserRepository)
	cacheManager := cache.NewCacheManager(nil)
	s3Manager := &storage.S3Manager{}

	service := NewVendorService(vendorRepo, userRepo, s3Manager, cacheManager)

	tests := []struct {
		name     string
		storeName string
		desc     string
	}{
		{"", "", "empty name"},
		{"AB", "AB", "too short"},
	}

	for _, tt := range tests {
		vendor := &models.Vendor{
			UserID:    1,
			StoreName: tt.storeName,
		}

		err := service.ValidateVendor(vendor)
		assert.Error(t, err, tt.desc)
	}
}

func TestValidateVendor_MissingSlug(t *testing.T) {
	vendorRepo := new(MockVendorRepository)
	userRepo := new(MockUserRepository)
	cacheManager := cache.NewCacheManager(nil)
	s3Manager := &storage.S3Manager{}

	service := NewVendorService(vendorRepo, userRepo, s3Manager, cacheManager)

	vendor := &models.Vendor{
		UserID:    1,
		StoreName: "Test Store",
		StoreSlug: "",
	}

	err := service.ValidateVendor(vendor)
	assert.Error(t, err)
}

func TestValidateVendor_InvalidCommissionModel(t *testing.T) {
	vendorRepo := new(MockVendorRepository)
	userRepo := new(MockUserRepository)
	cacheManager := cache.NewCacheManager(nil)
	s3Manager := &storage.S3Manager{}

	service := NewVendorService(vendorRepo, userRepo, s3Manager, cacheManager)

	vendor := &models.Vendor{
		UserID:          1,
		StoreName:       "Test Store",
		StoreSlug:       "test-store",
		CommissionModel: "invalid",
	}

	err := service.ValidateVendor(vendor)
	assert.Error(t, err)
}

func TestValidateVendor_InvalidCommissionRate_Negative(t *testing.T) {
	vendorRepo := new(MockVendorRepository)
	userRepo := new(MockUserRepository)
	cacheManager := cache.NewCacheManager(nil)
	s3Manager := &storage.S3Manager{}

	service := NewVendorService(vendorRepo, userRepo, s3Manager, cacheManager)

	vendor := &models.Vendor{
		UserID:          1,
		StoreName:       "Test Store",
		StoreSlug:       "test-store",
		CommissionModel: "margin",
		CommissionRate:  -0.1,
	}

	err := service.ValidateVendor(vendor)
	assert.Error(t, err)
}

func TestValidateVendor_InvalidCommissionRate_TooHigh(t *testing.T) {
	vendorRepo := new(MockVendorRepository)
	userRepo := new(MockUserRepository)
	cacheManager := cache.NewCacheManager(nil)
	s3Manager := &storage.S3Manager{}

	service := NewVendorService(vendorRepo, userRepo, s3Manager, cacheManager)

	vendor := &models.Vendor{
		UserID:          1,
		StoreName:       "Test Store",
		StoreSlug:       "test-store",
		CommissionModel: "markup",
		CommissionRate:  1.5,
	}

	err := service.ValidateVendor(vendor)
	assert.Error(t, err)
}

// Test UpdateBankDetails

func TestUpdateBankDetails_Success(t *testing.T) {
	vendorRepo := new(MockVendorRepository)
	userRepo := new(MockUserRepository)
	cacheManager := cache.NewCacheManager(nil)
	s3Manager := &storage.S3Manager{}

	service := NewVendorService(vendorRepo, userRepo, s3Manager, cacheManager)

	vendorRepo.On("UpdateBankDetails", mock.MatchedBy(func(details *models.VendorBankDetails) bool {
		return details.VendorID == 1 && details.AccountName == "John Doe"
	})).Return(nil)

	err := service.UpdateBankDetails(1, "John Doe", "1234567890", "Bank Name", "Branch")

	require.NoError(t, err)
	vendorRepo.AssertExpectations(t)
}

func TestUpdateBankDetails_Failed(t *testing.T) {
	vendorRepo := new(MockVendorRepository)
	userRepo := new(MockUserRepository)
	cacheManager := cache.NewCacheManager(nil)
	s3Manager := &storage.S3Manager{}

	service := NewVendorService(vendorRepo, userRepo, s3Manager, cacheManager)

	vendorRepo.On("UpdateBankDetails", mock.Anything).Return(errors.New("database error"))

	err := service.UpdateBankDetails(1, "John Doe", "1234567890", "Bank Name", "Branch")

	require.Error(t, err)
}

// Test GetVendorWallet

func TestGetVendorWallet_Success(t *testing.T) {
	vendorRepo := new(MockVendorRepository)
	userRepo := new(MockUserRepository)
	cacheManager := cache.NewCacheManager(nil)
	s3Manager := &storage.S3Manager{}

	service := NewVendorService(vendorRepo, userRepo, s3Manager, cacheManager)

	wallet := &models.VendorWallet{
		ID:       1,
		VendorID: 1,
		Balance:  5000.0,
	}

	vendorRepo.On("GetWallet", uint(1)).Return(wallet, nil)

	result, err := service.GetVendorWallet(1)

	require.NoError(t, err)
	assert.Equal(t, wallet, result)
}

func TestGetVendorWallet_NotFound(t *testing.T) {
	vendorRepo := new(MockVendorRepository)
	userRepo := new(MockUserRepository)
	cacheManager := cache.NewCacheManager(nil)
	s3Manager := &storage.S3Manager{}

	service := NewVendorService(vendorRepo, userRepo, s3Manager, cacheManager)

	vendorRepo.On("GetWallet", uint(999)).Return(nil, errors.New("wallet not found"))

	result, err := service.GetVendorWallet(999)

	require.Error(t, err)
	assert.Nil(t, result)
}
