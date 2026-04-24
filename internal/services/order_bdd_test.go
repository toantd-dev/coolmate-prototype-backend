// BDD-style acceptance tests for the Order service.
//
// Covers stories:
//   US-INV-001 · Stock is validated against cart quantity before checkout
//   US-CHK-001 · SplitOrder fans a master order into per-vendor SubOrders
//   US-CHK-002 · Promotions reduce the order total
//
// Reuses MockOrderRepository, MockPromotionRepository, MockVendorRepository,
// MockProductRepository from existing *_test.go files.
package services

import (
	"errors"
	"testing"
	"time"

	"github.com/coolmate/ecommerce-backend/internal/models"
	"github.com/coolmate/ecommerce-backend/pkg/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newOrderSvcBDD(t *testing.T) (
	*OrderService,
	*MockOrderRepository,
	*MockProductRepository,
	*MockPromotionRepository,
	*MockVendorRepository,
) {
	t.Helper()
	or := new(MockOrderRepository)
	pr := new(MockProductRepository)
	promoRepo := new(MockPromotionRepository)
	vr := new(MockVendorRepository)
	cm := cache.NewCacheManager(nil)
	return NewOrderService(or, pr, promoRepo, vr, cm), or, pr, promoRepo, vr
}

// =============================================================================
// US-INV-001 · ValidateStock refuses checkout if any variant is short
//
// AC: US-INV-001 AC1, AC2
// =============================================================================

func TestUS_INV_001_ValidateStock(t *testing.T) {
	// AC1
	// Given  every cart item has Variant.Stock ≥ Quantity
	// When   ValidateStock is called
	// Then   returns nil
	t.Run("AC1/sufficient_stock_passes", func(t *testing.T) {
		svc, _, _, _, _ := newOrderSvcBDD(t)
		items := []models.CartItem{
			{ID: 1, Quantity: 2, Variant: &models.ProductVariant{Stock: 10}},
			{ID: 2, Quantity: 5, Variant: &models.ProductVariant{Stock: 5}},
		}

		err := svc.ValidateStock(items)

		require.NoError(t, err, "every variant has enough stock")
	})

	// AC2
	// Given  at least one item where Variant.Stock < Quantity
	// When   ValidateStock is called
	// Then   an insufficient-stock error identifying the variant
	t.Run("AC2/insufficient_stock_errors", func(t *testing.T) {
		svc, _, _, _, _ := newOrderSvcBDD(t)
		items := []models.CartItem{
			{ID: 1, Quantity: 2, Variant: &models.ProductVariant{Stock: 10}},
			{ID: 2, Quantity: 6, Variant: &models.ProductVariant{Stock: 5}}, // short
		}

		err := svc.ValidateStock(items)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "stock", "message should mention stock")
	})
}

// =============================================================================
// US-CHK-001 · SplitOrder fans master order into per-vendor SubOrders
//
// AC: US-CHK-001 AC1, AC2, AC3, AC4
// =============================================================================

func TestUS_CHK_001_SplitOrder(t *testing.T) {
	// AC1
	// Given  3 cart items from vendor 10
	// When   SplitOrder is called
	// Then   one SubOrder with VendorID=10, subtotal = sum of line totals,
	//        VendorEarning = Subtotal − CommissionAmount
	t.Run("AC1/single_vendor_yields_one_suborder", func(t *testing.T) {
		svc, _, pr, _, vr := newOrderSvcBDD(t)
		vr.On("GetByID", uint(10)).Return(&models.Vendor{
			ID: 10, CommissionModel: "margin", CommissionRate: 0.05,
		}, nil)
		pr.On("GetCategory", mock.Anything).Return(nil, errors.New("no category"))

		order := &models.Order{ID: 1}
		items := []models.CartItem{
			{ID: 1, UnitPrice: 50, Quantity: 2, Product: &models.Product{ID: 1, VendorID: 10}},
			{ID: 2, UnitPrice: 80, Quantity: 1, Product: &models.Product{ID: 2, VendorID: 10}},
			{ID: 3, UnitPrice: 20, Quantity: 5, Product: &models.Product{ID: 3, VendorID: 10}},
		}
		// Expected subtotal = 100 + 80 + 100 = 280
		// margin @ 5% → commission = 14, vendor earning = 266

		subs, err := svc.SplitOrder(order, items)

		require.NoError(t, err)
		require.Len(t, subs, 1)
		s := subs[0]
		assert.Equal(t, uint(10), s.VendorID)
		assert.InDelta(t, 280.0, s.Subtotal, 1e-9)
		assert.InDelta(t, 14.0, s.CommissionAmount, 1e-9)
		assert.InDelta(t, 266.0, s.VendorEarning, 1e-9,
			"vendor earning = subtotal − commission")
	})

	// AC2
	// Given  items from vendor 10 and vendor 20
	// When   SplitOrder is called
	// Then   two SubOrders, each with its own subtotal and commission
	t.Run("AC2/multiple_vendors_yield_one_suborder_each", func(t *testing.T) {
		svc, _, pr, _, vr := newOrderSvcBDD(t)
		vr.On("GetByID", uint(10)).Return(&models.Vendor{
			ID: 10, CommissionModel: "margin", CommissionRate: 0.05,
		}, nil)
		vr.On("GetByID", uint(20)).Return(&models.Vendor{
			ID: 20, CommissionModel: "margin", CommissionRate: 0.10,
		}, nil)
		pr.On("GetCategory", mock.Anything).Return(nil, errors.New("no category"))

		order := &models.Order{ID: 1}
		items := []models.CartItem{
			{ID: 1, UnitPrice: 100, Quantity: 1, Product: &models.Product{ID: 1, VendorID: 10}}, // V10 subtotal 100
			{ID: 2, UnitPrice: 50, Quantity: 2, Product: &models.Product{ID: 2, VendorID: 10}},  // V10 subtotal 100 → total 200
			{ID: 3, UnitPrice: 60, Quantity: 1, Product: &models.Product{ID: 3, VendorID: 20}},  // V20 subtotal 60
		}

		subs, err := svc.SplitOrder(order, items)

		require.NoError(t, err)
		require.Len(t, subs, 2)

		byVendor := map[uint]*models.SubOrder{}
		for i := range subs {
			byVendor[subs[i].VendorID] = &subs[i]
		}
		assert.InDelta(t, 200.0, byVendor[10].Subtotal, 1e-9, "vendor 10 subtotal")
		assert.InDelta(t, 10.0, byVendor[10].CommissionAmount, 1e-9, "vendor 10 commission 5%")
		assert.InDelta(t, 60.0, byVendor[20].Subtotal, 1e-9, "vendor 20 subtotal")
		assert.InDelta(t, 6.0, byVendor[20].CommissionAmount, 1e-9, "vendor 20 commission 10%")
	})

	// AC3
	// Given  split succeeds
	// Then   order.Subtotal = sum of SubOrder.Subtotal  AND  GrandTotal reflects it
	//
	// NOTE: closes the "AC3 gap" flagged in USER_STORIES.md §5 CHK-001.
	t.Run("AC3/totals_rollup_on_master_order", func(t *testing.T) {
		svc, _, pr, _, vr := newOrderSvcBDD(t)
		vr.On("GetByID", uint(10)).Return(&models.Vendor{
			ID: 10, CommissionModel: "margin", CommissionRate: 0.05,
		}, nil)
		vr.On("GetByID", uint(20)).Return(&models.Vendor{
			ID: 20, CommissionModel: "margin", CommissionRate: 0.10,
		}, nil)
		pr.On("GetCategory", mock.Anything).Return(nil, errors.New("no category"))

		order := &models.Order{ID: 1}
		items := []models.CartItem{
			{ID: 1, UnitPrice: 100, Quantity: 2, Product: &models.Product{ID: 1, VendorID: 10}}, // 200
			{ID: 2, UnitPrice: 50, Quantity: 3, Product: &models.Product{ID: 2, VendorID: 20}},  // 150
		}

		subs, err := svc.SplitOrder(order, items)

		require.NoError(t, err)
		var sumSub float64
		for _, s := range subs {
			sumSub += s.Subtotal
		}
		// SplitOrder updates order.Subtotal; verify the invariant holds.
		assert.InDelta(t, sumSub, order.Subtotal, 1e-9,
			"order.Subtotal must equal sum of SubOrder subtotals")
		// GrandTotal is Subtotal + tax/shipping/etc; at minimum it must not
		// be less than Subtotal.
		assert.True(t, order.GrandTotal >= order.Subtotal,
			"GrandTotal (%v) must be ≥ Subtotal (%v)", order.GrandTotal, order.Subtotal)
	})

	// AC4
	// Given  invalid inputs: nil order, empty items, or item whose product/vendor
	//        does not load
	// When   SplitOrder is called
	// Then   a descriptive error; no SubOrders
	t.Run("AC4a/nil_order_errors", func(t *testing.T) {
		svc, _, _, _, _ := newOrderSvcBDD(t)
		subs, err := svc.SplitOrder(nil, []models.CartItem{{ID: 1}})
		require.Error(t, err)
		assert.Empty(t, subs)
	})

	t.Run("AC4b/empty_items_errors", func(t *testing.T) {
		svc, _, _, _, _ := newOrderSvcBDD(t)
		subs, err := svc.SplitOrder(&models.Order{ID: 1}, []models.CartItem{})
		require.Error(t, err)
		assert.Empty(t, subs)
	})

	t.Run("AC4c/missing_vendor_errors", func(t *testing.T) {
		svc, _, _, _, vr := newOrderSvcBDD(t)
		vr.On("GetByID", uint(999)).Return(nil, errors.New("vendor not found"))
		items := []models.CartItem{
			{ID: 1, UnitPrice: 10, Quantity: 1, Product: &models.Product{ID: 1, VendorID: 999}},
		}
		subs, err := svc.SplitOrder(&models.Order{ID: 1}, items)
		require.Error(t, err)
		assert.Empty(t, subs)
	})
}

// =============================================================================
// US-CHK-002 · Promotions reduce the order total
//
// AC: US-CHK-002 AC1, AC2, AC3, AC4
// =============================================================================

func TestUS_CHK_002_ApplyPromotions(t *testing.T) {
	// Helper: build a ready promotion. The model uses DiscountValue + ValidFrom/ValidTo.
	activePromo := func(code, discType string, value float64) *models.Promotion {
		return &models.Promotion{
			Code:          code,
			DiscountType:  discType,
			DiscountValue: value,
			IsActive:      true,
			ValidFrom:     time.Now().Add(-24 * time.Hour),
			ValidTo:       time.Now().Add(24 * time.Hour),
			Type:          "coupon",
			FundingType:   "platform",
		}
	}

	// AC1
	// Given  a valid percent promo (rate 10) and subtotal 1000
	// When   ApplyPromotions is called
	// Then   returned discount = 100
	t.Run("AC1/percent_discount", func(t *testing.T) {
		svc, _, _, promo, _ := newOrderSvcBDD(t)
		promo.On("GetByCode", "PERCENT10").Return(activePromo("PERCENT10", "percent", 10), nil)

		order := &models.Order{Subtotal: 1000}
		disc, err := svc.ApplyPromotions(order, []string{"PERCENT10"})

		require.NoError(t, err)
		assert.InDelta(t, 100.0, disc, 1e-9)
	})

	// AC2
	// Given  a valid flat promo (value 50)
	// When   ApplyPromotions is called
	// Then   returned discount = 50
	t.Run("AC2/flat_discount", func(t *testing.T) {
		svc, _, _, promo, _ := newOrderSvcBDD(t)
		promo.On("GetByCode", "FLAT50").Return(activePromo("FLAT50", "flat", 50), nil)

		order := &models.Order{Subtotal: 1000}
		disc, err := svc.ApplyPromotions(order, []string{"FLAT50"})

		require.NoError(t, err)
		assert.InDelta(t, 50.0, disc, 1e-9)
	})

	// AC3
	// Given  a code that does not exist (or is expired or inactive)
	// When   ApplyPromotions is called
	// Then   returned discount = 0 AND an explanatory error.
	//
	// SPEC-DRIFT NOTE: USER_STORIES.md says "contribution is 0 and the call
	// does not error". The implementation (and the existing
	// TestApplyPromotions_InvalidCode) aborts with an error. Either:
	//   (a) change the code to match spec (fail-soft / skip unknown codes), or
	//   (b) update USER_STORIES.md to "(0, <explanatory error>)".
	// Until resolved, tests assert actual behavior.
	t.Run("AC3/invalid_code_returns_error_and_zero_discount", func(t *testing.T) {
		svc, _, _, promo, _ := newOrderSvcBDD(t)
		promo.On("GetByCode", "BOGUS").Return(nil, errors.New("not found"))

		order := &models.Order{Subtotal: 1000}
		disc, err := svc.ApplyPromotions(order, []string{"BOGUS"})

		require.Error(t, err, "unknown promo codes currently abort the call")
		assert.InDelta(t, 0.0, disc, 1e-9)
	})

	// AC4
	// Given  two valid promos (percent 10% + flat 50)
	// When   ApplyPromotions is called with both
	// Then   the combined discount is returned (sequential application)
	t.Run("AC4/multiple_promos_stack", func(t *testing.T) {
		svc, _, _, promo, _ := newOrderSvcBDD(t)
		promo.On("GetByCode", "PERCENT10").Return(activePromo("PERCENT10", "percent", 10), nil)
		promo.On("GetByCode", "FLAT50").Return(activePromo("FLAT50", "flat", 50), nil)

		order := &models.Order{Subtotal: 1000}
		disc, err := svc.ApplyPromotions(order, []string{"PERCENT10", "FLAT50"})

		require.NoError(t, err)
		assert.Greater(t, disc, 100.0, "multiple promos should combine beyond 100")
	})
}
