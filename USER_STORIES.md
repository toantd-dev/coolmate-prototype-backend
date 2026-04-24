# User Stories (MVP) — BDD Acceptance Criteria

Single reference contract for MVP backend unit tests. Each AC = one unit test, named `TestUS_<CODE>_<NNN>_AC<n>_<shortScenario>`.

**Scope:** the 10 MVP items from [REQUIREMENTS_ANALYSIS.md](REQUIREMENTS_ANALYSIS.md) § Critical Path to MVP.

**Reading a story:** `US-<CODE>-<NNN> <status> · Title` where status is **✅ tested** / **🟡 implemented, no test** / **❌ not implemented**. Pre-ACs header names the service function + file. Post-ACs footer lists matching `_test.go` tests (omitted when none) and gap ACs.

**Conventions in ACs**
- Given/When/Then describe **business state**, not mock setup. Error-propagation and cache side-effects only appear when they carry business meaning.
- Handler-layer binding tests are consolidated in one appendix section ([Request Validation Boundary](#13-request-validation-boundary)), not scattered inside service stories.

---

## Contents

1. [Authentication](#1-authentication) — AUTH (3 stories, all 🟡)
2. [Vendor KYC & Lifecycle](#2-vendor-kyc--lifecycle) — VEN (4 stories, all ✅)
3. [Product CRUD & Approval](#3-product-crud--approval) — PRD (4 stories, all ✅)
4. [Variants & Stock](#4-variants--stock) — INV (1 story, 🟡)
5. [Checkout & SplitOrder](#5-checkout--splitorder) — CHK (4 stories, mixed)
6. [Manual Transfer Payment](#6-manual-transfer-payment) — PAY (1 story, ❌)
7. [Order Lifecycle](#7-order-lifecycle) — ORD (3 stories, ❌)
8. [Commission](#8-commission) — COM (3 stories, all ✅)
9. [Returns & Refunds](#9-returns--refunds) — RTN (3 stories, ❌)
10. [Wallet & Settlement](#10-wallet--settlement) — WAL (3 stories, ❌)
11. [Vendor Notifications](#11-vendor-notifications) — NOT (2 stories, ❌)
12. [Request Validation Boundary](#12-request-validation-boundary) — VAL (1 story, ✅)

---

## 1. Authentication

### US-AUTH-001 🟡 · Register a new account
*As a new user, I want to register with email/password/name/role so I can access the platform with a token pair.*
**Traces:** `AuthService.Register` in [auth_service.go](internal/services/auth_service.go)

- **AC1** — **Given** no user with email E; **When** `Register(valid payload, role=customer)`; **Then** a User is created with a hashed password **And** access + refresh tokens are returned.
- **AC2** — **Given** a user with email E already exists; **When** `Register` is called with E; **Then** a conflict error is returned and no user row is written.

**Gaps:** AC1, AC2

---

### US-AUTH-002 🟡 · Login issues a token pair
*As a registered user, I want to exchange credentials for tokens so my API calls are authenticated.*
**Traces:** `AuthService.Login` in [auth_service.go](internal/services/auth_service.go)

- **AC1** — **Given** a user with matching password hash; **When** `Login(email, password)`; **Then** a 15-min access token + 7-day refresh token are returned **And** the refresh-token hash is cached in Redis under the user key.
- **AC2** — **Given** wrong password **or** unknown email; **When** `Login` is called; **Then** an unauthorized error is returned (identical shape to prevent user enumeration) **And** nothing is cached.

**Gaps:** AC1, AC2

---

### US-AUTH-003 🟡 · Refresh rotates, logout revokes
*As an authenticated client, I want refresh to extend the session and logout to revoke it so a leaked token can be killed.*
**Traces:** `AuthService.Refresh`, `AuthService.Logout` in [auth_service.go](internal/services/auth_service.go)

- **AC1 — Refresh rotates** — **Given** a valid refresh token cached for user 42; **When** `Refresh(token)`; **Then** a new access+refresh pair is returned **And** the old hash is replaced in Redis (rotation).
- **AC2 — Refresh rejects unknown / revoked token** — **Given** no cache entry matches the token; **When** `Refresh(token)`; **Then** unauthorized is returned.
- **AC3 — Logout revokes** — **Given** a cached refresh token; **When** `Logout(token)`; **Then** the Redis entry is deleted **And** a subsequent `Refresh` with the same token returns unauthorized.

**Gaps:** AC1, AC2, AC3

---

## 2. Vendor KYC & Lifecycle

### US-VEN-001 ✅ · Vendor registration validates store inputs
*As the platform, I want vendor input to be rejected before DB write so invalid vendors never exist.*
**Traces:** `VendorService.ValidateVendor` in [vendor_service.go](internal/services/vendor_service.go)

- **AC1 — Happy path** — **Given** `UserID>0`, `StoreName≥3`, StoreSlug set, CommissionModel `margin`, Rate `0.15`; **When** `ValidateVendor(v)`; **Then** nil.
- **AC2 — Model must be margin/markup** — **Given** `CommissionModel="flat"`; **When** validate; **Then** error.
- **AC3 — Rate must be in [0,1]** — **Given** Rate = −0.1 or 1.5; **When** validate; **Then** error.
- **AC4 — Required fields** — **Given** UserID=0 **or** StoreSlug="" **or** StoreName<3 chars; **When** validate; **Then** error names the missing field.

**Tests:** `TestValidateVendor_Valid`, `_NilVendor`, `_MissingUserID`, `_InvalidStoreName`, `_MissingSlug`, `_InvalidCommissionModel`, `_InvalidCommissionRate_Negative`, `_InvalidCommissionRate_TooHigh`

---

### US-VEN-002 ✅ · Admin drives the vendor status machine
*As an admin, I want to approve, reject, or suspend vendors so the catalog reflects KYC decisions.*
**Traces:** `VendorService.ApproveVendor`, `RejectVendor`, `SuspendVendor` in [vendor_service.go](internal/services/vendor_service.go)

- **AC1 — Approve** — **Given** vendor status `pending`; **When** `ApproveVendor(id)`; **Then** status becomes `approved`.
- **AC2 — Reject** — **Given** vendor status `pending`; **When** `RejectVendor(id)`; **Then** status becomes `rejected`.
- **AC3 — Suspend** — **Given** vendor status `approved`; **When** `SuspendVendor(id)`; **Then** status becomes `suspended`.

**Tests:** `TestApproveVendor_Success`, `TestRejectVendor_Success`, `TestSuspendVendor_Success` (plus `TestApproveVendor_UpdateFailed` for repo-error propagation)

---

### US-VEN-003 ✅ · Product listing requires KYC approval + latest agreement
*As the platform, I want to gate product listing on both approval and agreement acceptance so only compliant vendors appear.*
**Traces:** `VendorService.CanListProducts` in [vendor_service.go](internal/services/vendor_service.go)

- **AC1 — Eligible** — **Given** `Status=approved` AND `AgreementAcceptedAt != nil`; **When** `CanListProducts(id)`; **Then** `(true, nil)`.
- **AC2 — Not approved** — **Given** `Status=pending|rejected|suspended`; **When** called; **Then** `(false, nil)`.
- **AC3 — Missing agreement** — **Given** `Status=approved` AND `AgreementAcceptedAt == nil`; **When** called; **Then** `(false, nil)`.
- **AC4 — Vendor not found** — **Given** no vendor with id; **When** called; **Then** `(false, <not-found error>)`.

**Tests:** `TestCanListProducts_Approved`, `_NotApproved`, `_NoAgreement`, `_VendorNotFound`

---

### US-VEN-004 ✅ · Admin updates vendor bank details
*As an admin, I want to set or change the vendor's bank account so settlements can be paid out.*
**Traces:** `VendorService.UpdateBankDetails` in [vendor_service.go](internal/services/vendor_service.go)

- **AC1 — Valid update persists** — **Given** a vendor exists; **When** `UpdateBankDetails(id, name≥5, number≥8, bank≥2, branch≥2)`; **Then** the new details are persisted.

**Tests:** `TestUpdateBankDetails_Success`, `TestUpdateBankDetails_Failed`

---

## 3. Product CRUD & Approval

### US-PRD-001 ✅ · Vendor creates a draft product
*As a vendor, I want new products to land in `draft` so I can prepare them before admin review.*
**Traces:** `ProductService.CreateProduct` in [product_service.go](internal/services/product_service.go)

- **AC1 — Valid product created as draft** — **Given** Name 3–255, VendorID>0, CategoryID>0, BasePrice>0, CostPrice<BasePrice; **When** `CreateProduct(p)`; **Then** inserted with `Status="draft"`.
- **AC2 — Cost ≥ base rejected** — **Given** CostPrice = BasePrice or greater; **When** create; **Then** validation error, no insert.
- **AC3 — Returnable without window rejected** — **Given** `IsReturnable=true` AND `ReturnWindowDays≤0`; **When** create; **Then** validation error.

**Tests:** `TestCreateProduct_Success`, `TestCreateProduct_ValidationFailed`, `TestValidateProduct_*` suite

---

### US-PRD-002 ✅ · Updating a published product re-enters review
*As the platform, I want any edit to a published product to require re-approval so stale approvals cannot cover new content.*
**Traces:** `ProductService.UpdateProduct` in [product_service.go](internal/services/product_service.go)

- **AC1 — Published → pending_approval** — **Given** existing product `Status="published"`; **When** `UpdateProduct(p)`; **Then** status is set to `pending_approval` before save.
- **AC2 — Draft stays draft** — **Given** `Status="draft"`; **When** `UpdateProduct`; **Then** status stays `draft`.

**Tests:** `TestUpdateProduct_PublishedToApproval`, `TestUpdateProduct_ValidationFailed`

---

### US-PRD-003 ✅ · Admin decides on a pending product
*As an admin, I want to approve or reject pending products with a reason so vendors get clear publishing decisions.*
**Traces:** `ProductService.ApproveProduct`, `RejectProduct` in [product_service.go](internal/services/product_service.go)

- **AC1 — Approve publishes** — **Given** `Status="pending_approval"`; **When** `ApproveProduct(id)`; **Then** `Status="published"`.
- **AC2 — Reject stores reason** — **Given** `Status="pending_approval"`; **When** `RejectProduct(id, "blurry photos")`; **Then** `Status="rejected"` and the reason is persisted.

**Tests:** `TestApproveProduct_Success`, `TestRejectProduct_Success`

---

### US-PRD-004 ✅ · Product reads are cached
*As the platform, I want single-product reads to hit Redis before the DB so detail pages scale.*
**Traces:** `ProductService.GetProductByID` (5-min TTL), `GetCategories` (24-hr TTL) in [product_service.go](internal/services/product_service.go)

- **AC1 — Cache hit skips DB** — **Given** product 7 is cached; **When** `GetProductByID(7)`; **Then** the repo is not called and the cached value is returned.
- **AC2 — Cache miss populates cache** — **Given** no cache entry, DB has it; **When** `GetProductByID(7)`; **Then** the DB value is returned **And** written back to cache.
- **AC3 — Not found** — **Given** neither cache nor DB has product 99; **When** `GetProductByID(99)`; **Then** a not-found error is returned.

**Tests:** `TestGetProductByID_FromCache`, `_FromDatabase`, `_NotFound`, `TestGetCategories_FromCache`, `_FromDatabase`

---

## 4. Variants & Stock

### US-INV-001 🟡 · Stock is validated against cart quantity before checkout
*As the platform, I want to refuse checkout if any variant is short so we never oversell.*
**Traces:** `OrderService.ValidateStock` in [order_service.go](internal/services/order_service.go)

- **AC1 — Sufficient** — **Given** every item has `Variant.Stock ≥ Quantity`; **When** `ValidateStock(items)`; **Then** nil.
- **AC2 — Insufficient** — **Given** at least one item where `Variant.Stock < Quantity`; **When** validate; **Then** an insufficient-stock error identifying the variant.

**Tests:** `TestValidateStock_SufficientStock`, `TestValidateStock_InsufficientStock`

*Variant uniqueness (duplicate-SKU conflict) is not yet a service responsibility — tracked as an implementation gap, not a story here.*

---

## 5. Checkout & SplitOrder

### US-CHK-001 ✅ · SplitOrder fans a master order into per-vendor SubOrders
*As the checkout flow, I want one cart to produce one master Order plus one SubOrder per vendor with commission calculated so vendors settle independently.*
**Traces:** `OrderService.SplitOrder` in [order_service.go](internal/services/order_service.go)

- **AC1 — Single vendor** — **Given** 3 cart items from vendor 10; **When** `SplitOrder(order, items)`; **Then** one SubOrder with `VendorID=10`, Subtotal = sum of line totals, `VendorEarning = Subtotal − CommissionAmount`.
- **AC2 — Multiple vendors** — **Given** 2 items from vendor 10 and 1 from vendor 20; **When** split; **Then** two SubOrders, each with its own subtotal and commission.
- **AC3 — Totals roll up** — **Given** split succeeds; **Then** `order.Subtotal` = sum of SubOrder.Subtotal and `order.GrandTotal` reflects it.
- **AC4 — Invalid inputs** — **Given** `order==nil` **or** `items==[]` **or** an item whose product/vendor does not load; **When** split; **Then** a descriptive error and no SubOrders are persisted.

**Tests:** `TestSplitOrder_SingleVendor`, `_MultipleVendors`, `_NilOrder`, `_EmptyCartItems`, `_CartItemWithoutProduct`, `_VendorNotFound`
**Gaps:** AC3 (verify totals are asserted; otherwise add)

---

### US-CHK-002 ✅ · Promotions reduce the order total
*As a customer, I want valid promo codes to discount my order so advertised offers are honoured.*
**Traces:** `OrderService.ApplyPromotions` in [order_service.go](internal/services/order_service.go)

- **AC1 — Percent discount** — **Given** a valid promo `percent, 10` and subtotal 1000; **When** `ApplyPromotions(order, [code])`; **Then** returned discount = 100.
- **AC2 — Flat discount** — **Given** valid promo `flat, 50`; **When** apply; **Then** returned discount = 50.
- **AC3 — Invalid code is skipped** — **Given** a code that does not exist, is expired, or is inactive; **When** apply; **Then** its contribution is 0 and the call does not error.
- **AC4 — Multiple promos stack** — **Given** two valid promos; **When** apply both codes; **Then** the combined discount is returned.

**Tests:** `TestApplyPromotions_ValidPromo`, `_InvalidCode`, `_FlatDiscount`, `_MultiplePromos`

---

### US-CHK-003 ❌ · Checkout rolls back stock on failure
*As the platform, I want a failed checkout to leave inventory untouched so we don't leak stock on errors.*
**Traces:** not yet implemented — no orchestrating `Checkout` function composes stock decrement + split.

- **AC1 — Rollback on split failure** — **Given** stock was decremented for all items; **When** `SplitOrder` then errors (e.g., vendor not found); **Then** every decremented variant is restored AND no Order/SubOrder rows remain.
- **AC2 — Rollback on promotion failure** — **Given** stock was decremented; **When** `ApplyPromotions` errors; **Then** stock is restored.

**Gaps:** AC1, AC2 (implementation + test both needed)

---

### US-CHK-004 ❌ · Checkout creates a pending wallet credit per vendor
*As the platform, I want each vendor to accrue a pending wallet credit on successful checkout so earnings are trackable before settlement.*
**Traces:** not yet implemented — `VendorWallet` exists, no `RecordTransaction` service.

- **AC1 — One pending credit per SubOrder** — **Given** split produces N SubOrders; **When** checkout commits; **Then** N `WalletTransaction` rows are created, each `type=credit, status=pending, amount = SubOrder.VendorEarning`.

**Gaps:** AC1

---

## 6. Manual Transfer Payment

### US-PAY-001 ❌ · Manual-transfer payment lifecycle
*As a customer paying by bank transfer, I want to submit proof; as an admin, I want to confirm it so the order can proceed.*
**Traces:** not yet implemented (`ProcessManualTransfer`, `MarkPaymentPaid` missing). Validator already accepts `manual_transfer` enum.

- **AC1 — Customer submits proof** — **Given** a pending order with `PaymentMethod="manual_transfer"`; **When** the customer uploads proof (file/URL); **Then** a payment-proof record is linked to the order **And** `PaymentStatus` stays `pending`.
- **AC2 — Admin cannot mark paid without proof** — **Given** no proof exists; **When** `MarkPaymentPaid(orderID)`; **Then** a precondition error.
- **AC3 — Admin confirms payment** — **Given** proof is uploaded; **When** admin calls `MarkPaymentPaid`; **Then** `PaymentStatus="paid"` **And** vendor new-order notification fires.

**Gaps:** AC1, AC2, AC3

---

## 7. Order Lifecycle

### US-ORD-001 ❌ · Vendor drives SubOrder fulfillment
*As a vendor, I want to mark my SubOrder ready-to-ship and then shipped (with tracking) so the customer knows the parcel moves.*
**Traces:** not yet implemented as a guarded transition; repository-level updates exist.

- **AC1 — Valid pending → ready_to_ship** — **Given** SubOrder `pending` and caller is its vendor; **When** transition to `ready_to_ship` requested; **Then** status updates.
- **AC2 — Tracking required for shipped** — **Given** `ready_to_ship`; **When** transition to `shipped` without tracking number; **Then** a validation error.
- **AC3 — Shipped with tracking** — **Given** tracking number provided; **When** transition; **Then** `shipped`, tracking stored, customer shipment notification fires.
- **AC4 — Forbidden skips and cross-vendor edits** — **Given** vendor tries to skip to `delivered` **OR** caller's VendorID does not own the SubOrder; **When** transition; **Then** forbidden/unauthorized error.

**Gaps:** AC1, AC2, AC3, AC4

---

### US-ORD-002 ❌ · Delivered master order finalizes commissions
*As an admin, I want the master order to reach `delivered` only when all SubOrders are delivered, and for vendor wallets to move from pending to available at that moment.*
**Traces:** not yet implemented.

- **AC1 — All SubOrders delivered → master delivered** — **Given** every SubOrder of order 42 is `delivered`; **When** the last one is marked; **Then** `order.Status="delivered"`.
- **AC2 — Wallet transactions finalize** — **Given** master order transitions to `delivered`; **When** the hook runs; **Then** each vendor's matching pending `WalletTransaction` moves to `available`.

**Gaps:** AC1, AC2

---

### US-ORD-003 ❌ · Cancel reverts stock and wallet
*As a customer (before processing) or admin (any time), I want cancellation to undo inventory and earnings side-effects so records stay consistent.*
**Traces:** not yet implemented as a service orchestration.

- **AC1 — Cancel before processing** — **Given** order `pending`; **When** customer cancels; **Then** status `cancelled`, all stock decrements reverted, pending wallet transactions deleted.
- **AC2 — Cancel after payment triggers refund** — **Given** order `paid`; **When** admin cancels; **Then** status `cancelled`, a Refund is created, vendor earnings are reversed.

**Gaps:** AC1, AC2

---

## 8. Commission

### US-COM-001 ✅ · Commission hierarchy: category → vendor → platform default
*As the platform, I want category commission to override vendor commission, and a 5% margin default to apply when neither is set, so category configuration always wins.*
**Traces:** `CommissionService.CalculateCommission` in [commission_service.go](internal/services/commission_service.go)

- **AC1 — Category wins** — **Given** category `(margin, 0.20)` AND vendor `(markup, 0.30)`; **When** `CalculateCommission(item, "markup", 0.30)`; **Then** the applied config is `(margin, 0.20)`.
- **AC2 — Vendor when no category** — **Given** category unset AND vendor `(margin, 0.15)`; **When** calculate; **Then** vendor rate/model is used.
- **AC3 — Platform default** — **Given** category unset AND `vendor.Model==""`; **When** calculate; **Then** `(margin, 0.05)` is used.
- **AC4 — Invalid inputs** — **Given** `orderItem==nil` **or** product missing; **When** calculate; **Then** an error is returned.

**Tests:** `TestCalculateCommission_CategoryCommission`, `_VendorCommissionWhenNoCategoryCommission`, `_PlatformDefaultWhenNoVendorCommission`, `_NilOrderItem`, `_NoProduct`, `TestGetCategoryCommission_Success`, `_NotFound`

---

### US-COM-002 ✅ · Margin vs markup formulas
*As the platform, I want margin and markup to produce the spec-documented values so finance numbers match.*
**Traces:** `calculateByModel` via `CalculateCommission` in [commission_service.go](internal/services/commission_service.go)

- **AC1 — Margin** — **Given** subtotal=100, rate=0.20; **Then** commission = 20 (`subtotal × rate`).
- **AC2 — Markup** — **Given** subtotal=100, rate=0.10; **Then** commission ≈ 9.09 (`subtotal / (1+rate) × (1 − 1/(1+rate))`).
- **AC3 — Unknown model** — **Given** `model="flat"`; **Then** commission = 0 (documented fallback).
- **AC4 — Zero edge cases** — **Given** quantity=0 **or** unit price=0; **Then** commission = 0 without error.

**Tests:** `TestCalculateByModel_Margin`, `_Markup`, `_UnknownModel`, `TestCalculateCommission_ZeroQuantity`, `_ZeroUnitPrice`, `_LargePrices`

---

### US-COM-003 ✅ · Commission rate must be in [0, 1]
*As the platform, I want rates outside [0,1] rejected so nonsensical payouts are impossible.*
**Traces:** `CommissionService.ValidateCommissionRate` in [commission_service.go](internal/services/commission_service.go)

- **AC1 — In range passes** — **Given** rate ∈ {0, 0.15, 1.0} with a known model; **Then** nil.
- **AC2 — Out of range fails** — **Given** rate = −0.01 or 1.01; **Then** a validation error.

**Tests:** `TestValidateCommissionRate_Valid`, `_Invalid`

---

## 9. Returns & Refunds

### US-RTN-001 ❌ · Customer initiates a return within the window
*As a customer, I want to open a return request while the product is still returnable so I can recover money on unsatisfactory items.*
**Traces:** not yet implemented (`InitiateReturn`).

- **AC1 — Within window + returnable** — **Given** order delivered D days ago, `product.IsReturnable=true`, `product.ReturnWindowDays ≥ D`; **When** `InitiateReturn(reason 10–1000 chars, ≤5 evidence URLs)`; **Then** a ReturnRequest row is created with `Status="requested"`.
- **AC2 — Window expired** — **Given** `D > product.ReturnWindowDays`; **When** initiate; **Then** precondition error.
- **AC3 — Not returnable** — **Given** `product.IsReturnable=false`; **When** initiate; **Then** precondition error.
- **AC4 — Duplicate blocked** — **Given** a return already exists for this order item; **When** initiate again; **Then** conflict error.

**Gaps:** AC1, AC2, AC3, AC4

---

### US-RTN-002 ❌ · Vendor or admin approves / rejects a return
*As a vendor or admin, I want to make an approval decision with a note so the refund pipeline can run (or not).*
**Traces:** not yet implemented.

- **AC1 — Approve** — **Given** `Status="requested"`; **When** `ApproveReturn(id, note)`; **Then** `Status="approved"`.
- **AC2 — Reject blocks refund** — **Given** `Status="requested"`; **When** `RejectReturn(id, reason)`; **Then** `Status="rejected"` and no Refund row is created.

**Gaps:** AC1, AC2

---

### US-RTN-003 ❌ · Processed refund reverses vendor earnings
*As the platform, I want a processed refund to debit the vendor wallet and increase the order's refund total so financial records stay balanced.*
**Traces:** not yet implemented (`ProcessRefund`, `ReverseVendorEarnings`).

- **AC1 — Standard refund** — **Given** a return `completed` with refund amount = vendor earning; **When** `ProcessRefund(returnID)`; **Then** a Refund row is `processed` **And** a debit `WalletTransaction` equal to the vendor share is appended **And** `order.refund_total` increases by the refund amount.
- **AC2 — Platform-funded discount not clawed back from vendor** — **Given** the original order had a platform-funded discount; **When** refund processes; **Then** the platform portion is not deducted from vendor earnings.

**Gaps:** AC1, AC2

---

## 10. Wallet & Settlement

### US-WAL-001 ❌ · Wallet ledger is append-only
*As the platform, I want every earning, reversal, and payout to be a new immutable row so the ledger is auditable.*
**Traces:** not yet implemented. Model: `models.WalletTransaction`.

- **AC1 — Credit on earning event** — **Given** a qualifying event (checkout, delivery); **When** `RecordTransaction(vendorID, {type:credit, amount, refType:order})`; **Then** a new row is inserted with those fields and a computed `balance_after`.
- **AC2 — Debit on refund / payout** — **Given** a qualifying event; **When** `RecordTransaction(…, type:debit)`; **Then** a new row is inserted and the running balance decreases.
- **AC3 — Immutable** — **Given** an existing row; **When** any update is attempted; **Then** it is refused.

**Gaps:** AC1, AC2, AC3

---

### US-WAL-002 ❌ · Admin previews a vendor's settlement for a period
*As an admin, I want a computed settlement preview so I can approve payouts with accurate numbers.*
**Traces:** not yet implemented (`GetSettlementCalculation`).

- **AC1 — Net payable formula** — **Given** delivered orders in `[start, end]` for vendor V; **When** `GetSettlementCalculation(V, start, end)`; **Then** returns `{grossSales, totalCommission, totalRefunds, netPayable}` where `netPayable = grossSales − totalCommission − totalRefunds`.
- **AC2 — Range is strict** — **Given** an order delivered outside the range; **When** calculate; **Then** it contributes to no total.

**Gaps:** AC1, AC2

---

### US-WAL-003 ❌ · Admin approves and pays a settlement
*As an admin, I want to approve and then mark a settlement paid with a reference so the wallet records the payout.*
**Traces:** not yet implemented (`ApproveSettlement`, `MarkSettlementAsPaid`).

- **AC1 — Approve** — **Given** `Status="pending"`; **When** `ApproveSettlement(id, note)`; **Then** `Status="processing"`.
- **AC2 — Mark paid debits wallet** — **Given** `Status="processing"`; **When** `MarkSettlementAsPaid(id, paymentRef)`; **Then** `Status="paid"`, `paid_at` is set, and a debit `WalletTransaction` equal to `netPayable` is appended.

**Gaps:** AC1, AC2

---

## 11. Vendor Notifications

### US-NOT-001 ❌ · Vendor receives a new-order email at checkout
*As a vendor, I want an email the moment checkout produces my SubOrder so I can ship quickly.*
**Traces:** not yet implemented (`SendNewOrderNotification`).

- **AC1 — One email per SubOrder vendor** — **Given** split produces SubOrders for vendors [10, 20]; **When** checkout commits; **Then** exactly one email is dispatched to vendor 10 and one to vendor 20.
- **AC2 — Email content** — **Given** the email is rendered; **Then** it contains order id, customer name, items, total, and shipping deadline.

**Gaps:** AC1, AC2

---

### US-NOT-002 ❌ · Vendor receives a status-change email on admin intervention
*As a vendor, I want an email when admin overrides a SubOrder's status so I don't miss interventions — but not when I'm the one making the change.*
**Traces:** not yet implemented (`SendOrderStatusChangeNotification`).

- **AC1 — Admin-initiated change emails the vendor** — **Given** admin transitions SubOrder 77 from `pending` to `cancelled`; **When** it commits; **Then** an email to the SubOrder's vendor containing `from`, `to`, reason.
- **AC2 — Vendor-initiated change does not self-email** — **Given** the vendor marks their own SubOrder `ready_to_ship`; **When** it commits; **Then** no email to that vendor.

**Gaps:** AC1, AC2

---

## 12. Request Validation Boundary

### US-VAL-001 ✅ · Handler binding rejects malformed payloads before the service layer
*As the platform, I want Gin's binding tags to catch malformed requests at the edge so invalid input never reaches business logic.*
**Traces:** struct tags in [request_validators.go](internal/handlers/request_validators.go); suite in [request_validators_test.go](internal/handlers/request_validators_test.go)

- **AC1 — Register rejects unsupported role / short password** — `role ∉ {vendor, customer}` OR `password len < 8` → binding error.
- **AC2 — Cart quantities are in [1, 1000]** — `AddToCartRequest.Quantity = 0` or `> 1000` → binding error.
- **AC3 — Coupon code length** — `ApplyCouponRequest.PromoCode == ""` or length > 50 → binding error.
- **AC4 — Checkout payment method is whitelisted** — `PaymentMethod ∉ {bank_ipg, emi, cod, manual_transfer}` → binding error.
- **AC5 — Checkout promo codes capped at 5** — `len(PromoCodes) > 5` → binding error.
- **AC6 — Shipping address shape** — Phone length ≠ 10 OR Country length ≠ 2 OR PostalCode length ≠ 5 → binding error.
- **AC7 — Return request shape** — `InitiateReturnRequest` reason length outside [10, 1000] or `len(Evidence) > 5` → binding error.

**Tests:** 35+ cases in `TestBind*` / `TestValidate*` of [request_validators_test.go](internal/handlers/request_validators_test.go). When writing new ACs for service stories, do **not** duplicate binding checks here — link to AC by number instead.

---

## Coverage Summary

| # | Module | Stories | ACs | Tested | Gap |
|---|---|--:|--:|--:|--:|
| 1 | Auth | 3 | 7 | 0 | 7 |
| 2 | Vendor KYC | 4 | 11 | 11 | 0 |
| 3 | Product | 4 | 10 | 10 | 0 |
| 4 | Variants & Stock | 1 | 2 | 2 | 0 |
| 5 | Checkout | 4 | 11 | 8 | 3 |
| 6 | Manual Payment | 1 | 3 | 0 | 3 |
| 7 | Order Lifecycle | 3 | 8 | 0 | 8 |
| 8 | Commission | 3 | 10 | 10 | 0 |
| 9 | Returns & Refunds | 3 | 8 | 0 | 8 |
| 10 | Wallet & Settlement | 3 | 7 | 0 | 7 |
| 11 | Vendor Notifications | 2 | 4 | 0 | 4 |
| 12 | Request Validation | 1 | 7 | 7 | 0 |
| | **Total** | **32** | **88** | **48** | **40** |

**Test-authoring priority (by gap impact):** Order Lifecycle (8), Returns & Refunds (8), Wallet & Settlement (7), Auth (7), Vendor Notifications (4), Manual Payment (3), Checkout rollback/wallet (3).
