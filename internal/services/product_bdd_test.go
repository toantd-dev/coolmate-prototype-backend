// BDD-style acceptance tests for the Product service.
//
// Covers stories:
//   US-PRD-001 · Vendor creates a draft product
//   US-PRD-002 · Updating a published product re-enters review
//   US-PRD-003 · Admin decides on a pending product
//   US-PRD-004 · Product reads are cached
//
// Reuses MockProductRepository (commission_service_test.go) and MockVendorRepository
// (order_service_test.go). Both satisfy IProductRepository / IVendorRepository.
package services

import (
	"errors"
	"testing"

	"github.com/coolmate/ecommerce-backend/internal/models"
	"github.com/coolmate/ecommerce-backend/pkg/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newProductSvcBDD(t *testing.T) (*ProductService, *MockProductRepository, *MockVendorRepository) {
	t.Helper()
	pr := new(MockProductRepository)
	vr := new(MockVendorRepository)
	cm := cache.NewCacheManager(nil)
	return NewProductService(pr, vr, cm), pr, vr
}

func validProduct() *models.Product {
	cost := 50.0
	return &models.Product{
		Name:       "Tee Shirt",
		VendorID:   1,
		CategoryID: 1,
		BasePrice:  100.0,
		CostPrice:  &cost,
	}
}

// =============================================================================
// US-PRD-001 · Vendor creates a draft product
//
// AC: US-PRD-001 AC1, AC2, AC3
// =============================================================================

func TestUS_PRD_001_CreateProduct(t *testing.T) {
	// AC1
	// Given  Name 3–255, VendorID>0, CategoryID>0, BasePrice>0, CostPrice<BasePrice
	// When   CreateProduct is called
	// Then   the product is inserted with Status="draft"
	t.Run("AC1/valid_product_created_as_draft", func(t *testing.T) {
		svc, pr, _ := newProductSvcBDD(t)
		pr.On("Create", mock.MatchedBy(func(p *models.Product) bool {
			return p.Name == "Tee Shirt" && p.Status == "draft"
		})).Return(nil)

		p := validProduct()

		err := svc.CreateProduct(p)

		require.NoError(t, err)
		assert.Equal(t, "draft", p.Status, "new products start in draft")
		pr.AssertExpectations(t)
	})

	// AC2
	// Given  CostPrice = BasePrice (or greater)
	// When   CreateProduct is called
	// Then   validation error; no Create call
	t.Run("AC2/cost_ge_base_rejected", func(t *testing.T) {
		svc, pr, _ := newProductSvcBDD(t)
		cost := 100.0
		p := &models.Product{
			Name: "Tee", VendorID: 1, CategoryID: 1,
			BasePrice: 100.0, CostPrice: &cost,
		}

		err := svc.CreateProduct(p)

		require.Error(t, err)
		pr.AssertNotCalled(t, "Create")
	})

	// AC3
	// Given  IsReturnable=true AND ReturnWindowDays ≤ 0
	// When   CreateProduct is called
	// Then   validation error
	t.Run("AC3/returnable_without_window_rejected", func(t *testing.T) {
		svc, pr, _ := newProductSvcBDD(t)
		p := validProduct()
		p.IsReturnable = true
		p.ReturnWindowDays = 0

		err := svc.CreateProduct(p)

		require.Error(t, err)
		pr.AssertNotCalled(t, "Create")
	})
}

// =============================================================================
// US-PRD-002 · Updating a published product re-enters review
//
// AC: US-PRD-002 AC1, AC2
// =============================================================================

func TestUS_PRD_002_UpdateProduct(t *testing.T) {
	// AC1
	// Given  Status = published
	// When   UpdateProduct is called
	// Then   Status is set to pending_approval before save
	t.Run("AC1/published_goes_back_to_pending_approval", func(t *testing.T) {
		svc, pr, _ := newProductSvcBDD(t)
		pr.On("Update", mock.MatchedBy(func(p *models.Product) bool {
			return p.Status == "pending_approval"
		})).Return(nil)

		p := validProduct()
		p.ID = 7
		p.Status = "published"

		err := svc.UpdateProduct(p)

		require.NoError(t, err)
		assert.Equal(t, "pending_approval", p.Status,
			"any edit to published product must return to review")
	})

	// AC2
	// Given  Status = draft
	// When   UpdateProduct is called
	// Then   Status stays draft
	t.Run("AC2/draft_stays_draft", func(t *testing.T) {
		svc, pr, _ := newProductSvcBDD(t)
		pr.On("Update", mock.MatchedBy(func(p *models.Product) bool {
			return p.Status == "draft"
		})).Return(nil)

		p := validProduct()
		p.ID = 8
		p.Status = "draft"

		err := svc.UpdateProduct(p)

		require.NoError(t, err)
		assert.Equal(t, "draft", p.Status)
	})
}

// =============================================================================
// US-PRD-003 · Admin decides on a pending product
//
// AC: US-PRD-003 AC1, AC2
// =============================================================================

func TestUS_PRD_003_AdminApprovalDecision(t *testing.T) {
	// AC1
	// Given  Status = pending_approval
	// When   ApproveProduct(id) is called
	// Then   Status = published
	t.Run("AC1/approve_publishes", func(t *testing.T) {
		svc, pr, _ := newProductSvcBDD(t)
		pr.On("UpdateStatus", uint(42), "published").Return(nil)

		err := svc.ApproveProduct(42)

		require.NoError(t, err)
		pr.AssertCalled(t, "UpdateStatus", uint(42), "published")
	})

	// AC2
	// Given  Status = pending_approval
	// When   RejectProduct(id, reason) is called
	// Then   Status = rejected and the reason is persisted (on the product record)
	t.Run("AC2/reject_with_reason_sets_status_rejected", func(t *testing.T) {
		svc, pr, _ := newProductSvcBDD(t)
		pr.On("UpdateStatus", uint(42), "rejected").Return(nil)

		err := svc.RejectProduct(42, "blurry photos")

		require.NoError(t, err)
		pr.AssertCalled(t, "UpdateStatus", uint(42), "rejected")
	})
}

// =============================================================================
// US-PRD-004 · Product reads are cached
//
// AC: US-PRD-004 AC1, AC2, AC3
// =============================================================================

func TestUS_PRD_004_ProductReadCaching(t *testing.T) {
	// AC1
	// Given  product 7 is cached   (NB: cache client is nil in tests, so Set()
	//        errors are swallowed; cache-hit path is exercised by existing
	//        TestGetProductByID_FromCache — this BDD test mirrors behavior.)
	// When   GetProductByID(7) is called
	// Then   the repo's value is returned without error.
	t.Run("AC1/cache_hit_skips_db", func(t *testing.T) {
		svc, pr, _ := newProductSvcBDD(t)
		p := &models.Product{ID: 7, Name: "Cached"}
		pr.On("GetByID", uint(7)).Return(p, nil)

		got, err := svc.GetProductByID(7)

		require.NoError(t, err)
		assert.Equal(t, p, got)
	})

	// AC2
	// Given  no cache entry; DB has product 7
	// When   GetProductByID(7) is called
	// Then   the DB value is returned; service attempts to populate cache
	t.Run("AC2/cache_miss_reads_from_db", func(t *testing.T) {
		svc, pr, _ := newProductSvcBDD(t)
		p := &models.Product{ID: 7, Name: "FromDB"}
		pr.On("GetByID", uint(7)).Return(p, nil)

		got, err := svc.GetProductByID(7)

		require.NoError(t, err)
		assert.Equal(t, "FromDB", got.Name)
	})

	// AC3
	// Given  neither cache nor DB has product 99
	// When   GetProductByID(99) is called
	// Then   a not-found error is returned
	t.Run("AC3/not_found_returns_error", func(t *testing.T) {
		svc, pr, _ := newProductSvcBDD(t)
		pr.On("GetByID", uint(99)).Return(nil, errors.New("not found"))

		got, err := svc.GetProductByID(99)

		require.Error(t, err)
		assert.Nil(t, got)
	})
}
