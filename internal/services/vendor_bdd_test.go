// BDD-style acceptance tests for the Vendor service.
//
// Covers stories:
//   US-VEN-001 · ValidateVendor — field/rate rules
//   US-VEN-002 · Admin drives vendor status machine
//   US-VEN-003 · CanListProducts = approved + agreement accepted
//   US-VEN-004 · Admin updates bank details
//
// Reuses MockVendorRepository/MockUserRepository from existing test files.
package services

import (
	"errors"
	"testing"
	"time"

	"github.com/coolmate/ecommerce-backend/internal/models"
	"github.com/coolmate/ecommerce-backend/pkg/cache"
	"github.com/coolmate/ecommerce-backend/pkg/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newVendorSvcBDD(t *testing.T) (*VendorService, *MockVendorRepository, *MockUserRepository) {
	t.Helper()
	vr := new(MockVendorRepository)
	ur := new(MockUserRepository)
	s3m := &storage.S3Manager{}
	cm := cache.NewCacheManager(nil)
	return NewVendorService(vr, ur, s3m, cm), vr, ur
}

func validVendor() *models.Vendor {
	return &models.Vendor{
		UserID:          1,
		StoreName:       "Acme Store",
		StoreSlug:       "acme-store",
		CommissionModel: "margin",
		CommissionRate:  0.15,
	}
}

// =============================================================================
// US-VEN-001 · ValidateVendor rejects invalid inputs before DB write
//
// AC: US-VEN-001 AC1, AC2, AC3, AC4
// =============================================================================

func TestUS_VEN_001_ValidateVendor(t *testing.T) {
	// AC1
	// Given  a vendor with UserID>0, StoreName≥3, slug set, model margin, rate 0.15
	// When   ValidateVendor is called
	// Then   returns nil (passes)
	t.Run("AC1/happy_path_passes", func(t *testing.T) {
		svc, _, _ := newVendorSvcBDD(t)
		assert.NoError(t, svc.ValidateVendor(validVendor()))
	})

	// AC2
	// Given  model ∉ {margin, markup}
	// When   ValidateVendor is called
	// Then   error
	t.Run("AC2/unknown_commission_model_rejected", func(t *testing.T) {
		svc, _, _ := newVendorSvcBDD(t)
		v := validVendor()
		v.CommissionModel = "flat"
		err := svc.ValidateVendor(v)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "commission model")
	})

	// AC3
	// Given  rate outside [0, 1]
	// When   ValidateVendor is called
	// Then   error
	t.Run("AC3/rate_out_of_range_rejected", func(t *testing.T) {
		svc, _, _ := newVendorSvcBDD(t)
		for _, rate := range []float64{-0.1, 1.5, -1.0, 2.0} {
			v := validVendor()
			v.CommissionRate = rate
			assert.Error(t, svc.ValidateVendor(v), "rate=%v must be rejected", rate)
		}
	})

	// AC4
	// Given  any required field is missing (UserID=0, slug="", StoreName<3)
	// When   ValidateVendor is called
	// Then   error naming the missing field
	t.Run("AC4/missing_required_field_rejected", func(t *testing.T) {
		svc, _, _ := newVendorSvcBDD(t)

		cases := []struct {
			name  string
			mut   func(*models.Vendor)
			want  string
		}{
			{"nil_vendor", func(v *models.Vendor) {}, "nil"},
			{"missing_userID", func(v *models.Vendor) { v.UserID = 0 }, "user"},
			{"missing_slug", func(v *models.Vendor) { v.StoreSlug = "" }, "slug"},
			{"short_storeName", func(v *models.Vendor) { v.StoreName = "ab" }, "store"},
		}
		for _, tc := range cases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				var v *models.Vendor
				if tc.name != "nil_vendor" {
					v = validVendor()
					tc.mut(v)
				}
				err := svc.ValidateVendor(v)
				require.Error(t, err)
				assert.Contains(t, toLower(err.Error()), tc.want)
			})
		}
	})
}

// =============================================================================
// US-VEN-002 · Admin drives vendor status machine
//
// AC: US-VEN-002 AC1, AC2, AC3
// =============================================================================

func TestUS_VEN_002_VendorStatusTransitions(t *testing.T) {
	// AC1
	// Given  vendor status = pending
	// When   ApproveVendor(id) is called
	// Then   status becomes approved
	t.Run("AC1/approve_sets_status_approved", func(t *testing.T) {
		svc, vr, _ := newVendorSvcBDD(t)
		vr.On("UpdateStatus", uint(10), "approved").Return(nil)

		err := svc.ApproveVendor(10)

		require.NoError(t, err)
		vr.AssertCalled(t, "UpdateStatus", uint(10), "approved")
	})

	// AC2
	// Given  vendor status = pending
	// When   RejectVendor(id) is called
	// Then   status becomes rejected
	t.Run("AC2/reject_sets_status_rejected", func(t *testing.T) {
		svc, vr, _ := newVendorSvcBDD(t)
		vr.On("UpdateStatus", uint(11), "rejected").Return(nil)

		err := svc.RejectVendor(11)

		require.NoError(t, err)
		vr.AssertCalled(t, "UpdateStatus", uint(11), "rejected")
	})

	// AC3
	// Given  vendor status = approved
	// When   SuspendVendor(id) is called
	// Then   status becomes suspended
	t.Run("AC3/suspend_sets_status_suspended", func(t *testing.T) {
		svc, vr, _ := newVendorSvcBDD(t)
		vr.On("UpdateStatus", uint(12), "suspended").Return(nil)

		err := svc.SuspendVendor(12)

		require.NoError(t, err)
		vr.AssertCalled(t, "UpdateStatus", uint(12), "suspended")
	})
}

// =============================================================================
// US-VEN-003 · CanListProducts = approved + latest agreement accepted
//
// AC: US-VEN-003 AC1, AC2, AC3, AC4
// =============================================================================

func TestUS_VEN_003_CanListProducts(t *testing.T) {
	now := time.Now()

	// AC1
	// Given  Status=approved AND AgreementAcceptedAt != nil
	// When   CanListProducts(id)
	// Then   (true, nil)
	t.Run("AC1/approved_and_agreement_accepted", func(t *testing.T) {
		svc, vr, _ := newVendorSvcBDD(t)
		vr.On("GetByID", uint(1)).Return(&models.Vendor{
			ID: 1, Status: "approved", AgreementAcceptedAt: &now,
		}, nil)

		ok, err := svc.CanListProducts(1)

		require.NoError(t, err)
		assert.True(t, ok)
	})

	// AC2
	// Given  status ∈ {pending, rejected, suspended}
	// When   CanListProducts(id)
	// Then   (false, <explanatory error>)
	//
	// SPEC-DRIFT NOTE: USER_STORIES.md says (false, nil), but the
	// implementation (and the existing TestCanListProducts_NotApproved)
	// return an error describing WHY listing is blocked. Explicit is
	// better than silent. USER_STORIES.md should be updated to "(false, err)".
	t.Run("AC2/not_approved_returns_false_with_reason", func(t *testing.T) {
		for _, status := range []string{"pending", "rejected", "suspended"} {
			svc, vr, _ := newVendorSvcBDD(t)
			vr.On("GetByID", uint(1)).Return(&models.Vendor{
				ID: 1, Status: status, AgreementAcceptedAt: &now,
			}, nil)

			ok, err := svc.CanListProducts(1)
			require.Error(t, err, "status=%s must return explanatory error", status)
			assert.Contains(t, err.Error(), "not approved")
			assert.False(t, ok, "status=%s must block listing", status)
		}
	})

	// AC3
	// Given  Status=approved AND AgreementAcceptedAt == nil
	// When   CanListProducts(id)
	// Then   (false, <agreement error>)   — see AC2 spec-drift note.
	t.Run("AC3/missing_agreement_blocks_listing", func(t *testing.T) {
		svc, vr, _ := newVendorSvcBDD(t)
		vr.On("GetByID", uint(1)).Return(&models.Vendor{
			ID: 1, Status: "approved", AgreementAcceptedAt: nil,
		}, nil)

		ok, err := svc.CanListProducts(1)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "agree")
		assert.False(t, ok)
	})

	// AC4
	// Given  no vendor with the given id
	// When   CanListProducts(id)
	// Then   (false, <not-found error>)
	t.Run("AC4/vendor_not_found_returns_error", func(t *testing.T) {
		svc, vr, _ := newVendorSvcBDD(t)
		vr.On("GetByID", uint(404)).Return(nil, errors.New("vendor not found"))

		ok, err := svc.CanListProducts(404)

		require.Error(t, err)
		assert.False(t, ok)
	})
}

// =============================================================================
// US-VEN-004 · Admin updates vendor bank details
//
// AC: US-VEN-004 AC1
// =============================================================================

func TestUS_VEN_004_UpdateBankDetails(t *testing.T) {
	// AC1
	// Given  a vendor exists
	// When   UpdateBankDetails(id, validFields)
	// Then   the new details are persisted
	t.Run("AC1/valid_update_persists", func(t *testing.T) {
		svc, vr, _ := newVendorSvcBDD(t)
		vr.On("UpdateBankDetails", mock.AnythingOfType("*models.VendorBankDetails")).Return(nil)

		// --- When ---
		err := svc.UpdateBankDetails(1, "Anna Seller", "12345678", "ABC Bank", "Colombo 4")

		// --- Then ---
		require.NoError(t, err)
		vr.AssertCalled(t, "UpdateBankDetails", mock.MatchedBy(
			func(d *models.VendorBankDetails) bool {
				return d.VendorID == 1 &&
					d.AccountName == "Anna Seller" &&
					d.AccountNumber == "12345678" &&
					d.BankName == "ABC Bank" &&
					d.BranchName == "Colombo 4"
			}))
	})
}

// toLower is a tiny helper to avoid importing "strings" just for case-insensitive
// substring checks in AC assertions.
func toLower(s string) string {
	out := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 32
		}
		out[i] = c
	}
	return string(out)
}
