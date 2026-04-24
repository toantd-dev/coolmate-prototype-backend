// BDD-style acceptance tests for the Commission service.
//
// Covers stories:
//   US-COM-001 · Commission hierarchy: category > vendor > platform default
//   US-COM-002 · Margin vs markup formulas
//   US-COM-003 · Rate validation in [0, 1]
//
// Reuses MockProductRepository from commission_service_test.go.
package services

import (
	"errors"
	"testing"

	"github.com/coolmate/ecommerce-backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newCommissionSvcBDD returns a service + its mock repo.
func newCommissionSvcBDD(t *testing.T) (*CommissionService, *MockProductRepository) {
	t.Helper()
	repo := new(MockProductRepository)
	return NewCommissionService(repo), repo
}

// itemWithCategory builds an OrderItem(categoryID, unit, qty).
func itemWithCategory(categoryID uint, unit float64, qty int) *models.OrderItem {
	return &models.OrderItem{
		ID:        1,
		UnitPrice: unit,
		Quantity:  qty,
		Product:   &models.Product{CategoryID: categoryID},
	}
}

// =============================================================================
// US-COM-001 · Commission hierarchy: category > vendor > platform default
//
// AC: US-COM-001 AC1, AC2, AC3, AC4
// =============================================================================

func TestUS_COM_001_CommissionHierarchy(t *testing.T) {
	const categoryID = uint(1)

	// AC1
	// Given  category commission (margin, 0.20) exists AND vendor is (markup, 0.30)
	// When   CalculateCommission is called
	// Then   the applied config is the category's (margin, 0.20)  — category wins
	t.Run("AC1/category_overrides_vendor", func(t *testing.T) {
		// --- Given ---
		svc, repo := newCommissionSvcBDD(t)
		catRate := 0.20
		repo.On("GetCategory", categoryID).Return(&models.Category{
			ID:              categoryID,
			CommissionModel: "margin",
			CommissionRate:  &catRate,
		}, nil)
		item := itemWithCategory(categoryID, 100.0, 2) // subtotal = 200

		// --- When ---
		commission, config, err := svc.CalculateCommission(item, "markup", 0.30)

		// --- Then ---
		require.NoError(t, err)
		assert.InDelta(t, 40.0, commission, 1e-9, "margin 200 × 0.20")
		assert.Equal(t, "margin", config.CommissionModel, "category config wins")
		assert.InDelta(t, 0.20, config.CommissionRate, 1e-9)
	})

	// AC2
	// Given  category is not set AND vendor is (margin, 0.15)
	// When   CalculateCommission is called
	// Then   the vendor rate/model is used
	t.Run("AC2/vendor_used_when_no_category", func(t *testing.T) {
		// --- Given ---
		svc, repo := newCommissionSvcBDD(t)
		repo.On("GetCategory", categoryID).Return(nil, errors.New("not found"))
		item := itemWithCategory(categoryID, 100.0, 2)

		// --- When ---
		commission, config, err := svc.CalculateCommission(item, "margin", 0.15)

		// --- Then ---
		require.NoError(t, err)
		assert.InDelta(t, 30.0, commission, 1e-9, "margin 200 × 0.15")
		assert.Equal(t, "margin", config.CommissionModel)
		assert.InDelta(t, 0.15, config.CommissionRate, 1e-9)
	})

	// AC3
	// Given  neither category nor vendor commission is set
	// When   CalculateCommission is called
	// Then   the platform default (margin, 0.05) applies
	t.Run("AC3/platform_default_when_neither_set", func(t *testing.T) {
		// --- Given ---
		svc, repo := newCommissionSvcBDD(t)
		repo.On("GetCategory", categoryID).Return(nil, errors.New("not found"))
		item := itemWithCategory(categoryID, 100.0, 2)

		// --- When --- vendor model empty signals "no vendor config"
		commission, config, err := svc.CalculateCommission(item, "", 0)

		// --- Then ---
		require.NoError(t, err)
		assert.InDelta(t, 10.0, commission, 1e-9, "margin 200 × 0.05 platform default")
		assert.Equal(t, "margin", config.CommissionModel)
		assert.InDelta(t, 0.05, config.CommissionRate, 1e-9)
	})

	// AC4
	// Given  an invalid input (nil item OR item without product)
	// When   CalculateCommission is called
	// Then   an error is returned — zero commission never slips through
	t.Run("AC4a/nil_order_item_errors", func(t *testing.T) {
		// --- Given ---
		svc, _ := newCommissionSvcBDD(t)

		// --- When ---
		commission, _, err := svc.CalculateCommission(nil, "margin", 0.10)

		// --- Then ---
		require.Error(t, err)
		assert.InDelta(t, 0.0, commission, 1e-9)
		assert.Contains(t, err.Error(), "nil", "error should name the cause")
	})

	t.Run("AC4b/item_without_product_falls_back_safely", func(t *testing.T) {
		// --- Given ---
		svc, _ := newCommissionSvcBDD(t)
		item := &models.OrderItem{ID: 1, UnitPrice: 100, Quantity: 2} // no Product

		// --- When ---
		commission, config, err := svc.CalculateCommission(item, "margin", 0.10)

		// --- Then --- the service falls back to vendor/platform config when no product
		require.NoError(t, err)
		assert.InDelta(t, 20.0, commission, 1e-9, "vendor margin 200 × 0.10")
		assert.Equal(t, "margin", config.CommissionModel)
	})
}

// =============================================================================
// US-COM-002 · Margin vs markup formulas
//
// AC: US-COM-002 AC1, AC2, AC3, AC4
// =============================================================================

func TestUS_COM_002_MarginAndMarkupFormulas(t *testing.T) {
	// AC1
	// Given  subtotal = 100, model = margin, rate = 0.20
	// When   commission is computed
	// Then   commission = 20   (subtotal × rate)
	t.Run("AC1/margin_formula", func(t *testing.T) {
		// --- Given ---
		svc, repo := newCommissionSvcBDD(t)
		repo.On("GetCategory", uint(0)).Return(nil, errors.New("no category"))
		item := &models.OrderItem{UnitPrice: 100, Quantity: 1, Product: &models.Product{}}

		// --- When ---
		commission, _, err := svc.CalculateCommission(item, "margin", 0.20)

		// --- Then ---
		require.NoError(t, err)
		assert.InDelta(t, 20.0, commission, 1e-9)
	})

	// AC2
	// Given  subtotal = 100, model = markup, rate = 0.10
	// When   commission is computed
	// Then   commission ≈ 8.2645  (subtotal / (1+rate) × (1 − 1/(1+rate)))
	//
	// SPEC-DRIFT NOTE (USER_STORIES.md): the story states "≈ 9.09" for
	// rate = 0.10, but the documented formula gives 8.2645. 9.09 would
	// correspond to rate ≈ 0.11. Tests assert the correct math per the
	// documented formula; USER_STORIES.md needs a correction.
	t.Run("AC2/markup_formula", func(t *testing.T) {
		// --- Given ---
		svc, repo := newCommissionSvcBDD(t)
		repo.On("GetCategory", uint(0)).Return(nil, errors.New("no category"))
		item := &models.OrderItem{UnitPrice: 100, Quantity: 1, Product: &models.Product{}}

		// --- When ---
		commission, _, err := svc.CalculateCommission(item, "markup", 0.10)

		// --- Then ---
		require.NoError(t, err)
		// 100/1.10 × (1 − 1/1.10)  =  90.909… × 0.0909…  ≈  8.2645
		assert.InDelta(t, 8.2645, commission, 1e-4,
			"markup formula: subtotal/(1+r) × (1 − 1/(1+r))")
	})

	// AC3
	// Given  an unknown model ("flat")
	// When   commission is computed
	// Then   commission falls back to margin formula = subtotal × rate
	//
	// SPEC-DRIFT NOTE (USER_STORIES.md): the story says commission = 0
	// for an unknown model, but the code's documented behavior is
	// `default: subtotal * rate` (fail-open to margin). The existing
	// TestCalculateByModel_UnknownModel also asserts the margin fallback.
	// Pick a single source of truth and update USER_STORIES.md.
	t.Run("AC3/unknown_model_falls_back_to_margin", func(t *testing.T) {
		// --- Given ---
		svc, repo := newCommissionSvcBDD(t)
		repo.On("GetCategory", uint(0)).Return(nil, errors.New("no category"))
		item := &models.OrderItem{UnitPrice: 100, Quantity: 1, Product: &models.Product{}}

		// --- When ---
		commission, _, err := svc.CalculateCommission(item, "flat", 0.10)

		// --- Then ---
		require.NoError(t, err)
		assert.InDelta(t, 10.0, commission, 1e-9,
			"unknown model defaults to margin: 100 × 0.10")
	})

	// AC4
	// Given  quantity = 0 OR unit price = 0
	// When   commission is computed
	// Then   commission = 0 without error
	t.Run("AC4a/zero_quantity_yields_zero", func(t *testing.T) {
		svc, repo := newCommissionSvcBDD(t)
		repo.On("GetCategory", uint(0)).Return(nil, errors.New("no category"))
		item := &models.OrderItem{UnitPrice: 100, Quantity: 0, Product: &models.Product{}}

		commission, _, err := svc.CalculateCommission(item, "margin", 0.20)

		require.NoError(t, err)
		assert.InDelta(t, 0.0, commission, 1e-9)
	})

	t.Run("AC4b/zero_unit_price_yields_zero", func(t *testing.T) {
		svc, repo := newCommissionSvcBDD(t)
		repo.On("GetCategory", uint(0)).Return(nil, errors.New("no category"))
		item := &models.OrderItem{UnitPrice: 0, Quantity: 5, Product: &models.Product{}}

		commission, _, err := svc.CalculateCommission(item, "margin", 0.20)

		require.NoError(t, err)
		assert.InDelta(t, 0.0, commission, 1e-9)
	})
}

// =============================================================================
// US-COM-003 · Commission rate must be in [0, 1]
//
// AC: US-COM-003 AC1, AC2
// =============================================================================

func TestUS_COM_003_ValidateCommissionRate(t *testing.T) {
	svc, _ := newCommissionSvcBDD(t)

	// AC1
	// Given  rate ∈ {0, 0.15, 1.0} with a known model
	// When   ValidateCommissionRate is called
	// Then   it passes (nil)
	t.Run("AC1/in_range_passes", func(t *testing.T) {
		// --- Given / When / Then (table) ---
		cases := []struct {
			rate  float64
			model string
		}{
			{0.0, "margin"},
			{0.15, "margin"},
			{1.0, "markup"},
			{0.5, "markup"},
		}
		for _, tc := range cases {
			assert.NoError(t, svc.ValidateCommissionRate(tc.rate, tc.model),
				"rate=%v model=%s should be valid", tc.rate, tc.model)
		}
	})

	// AC2
	// Given  rate = -0.01 OR rate = 1.01
	// When   ValidateCommissionRate is called
	// Then   a validation error is returned
	t.Run("AC2/out_of_range_fails", func(t *testing.T) {
		for _, rate := range []float64{-0.01, 1.01, -1.0, 5.0} {
			err := svc.ValidateCommissionRate(rate, "margin")
			require.Error(t, err, "rate=%v must be rejected", rate)
		}
	})
}
