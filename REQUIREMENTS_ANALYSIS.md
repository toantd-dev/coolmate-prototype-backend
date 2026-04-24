# Phân Tích Chi Tiết Các Functions Theo Nghiệp Vụ
**Dựa theo SRS Document V8 - 26 Jan 2026**

---

## 1. AUTHENTICATION & USER MANAGEMENT

### 1.1 Authentication Service Functions

**Login/Registration Flow:**
```
✓ Register(email, password, firstName, lastName, role)
  → Create user with role (customer/vendor)
  → Send email verification
  → Return access + refresh tokens
  
✓ Login(email, password)
  → Validate credentials
  → Generate JWT tokens (access: 15min, refresh: 7 days)
  → Cache refresh token in Redis
  → Return tokens
  
✓ VerifyEmail(token)
  → Verify email token
  → Mark email as verified
  → Enable product listing (for vendors)
  
✓ RefreshToken(refreshToken)
  → Validate refresh token
  → Generate new access token
  → Rotate refresh token
  
✓ Logout(userId)
  → Revoke refresh token (hash-based)
  → Clear Redis cache
  
✓ ForgotPassword(email)
  → Generate reset token (expires 1 hour)
  → Send email with reset link
  
✓ ResetPassword(token, newPassword)
  → Validate reset token
  → Hash new password
  → Update in database
```

**Social Login (Phase 2):**
```
✗ GoogleLogin(googleToken)
✗ FacebookLogin(facebookToken)
```

**OTP Verification (Phase 2):**
```
✗ SendOTP(phone)
✗ VerifyOTP(phone, otp)
```

---

## 2. VENDOR ONBOARDING & MANAGEMENT

### 2.1 Vendor Registration & KYC

**Registration Flow:**
```
✓ RegisterVendor(request)
  Fields:
  - Store name (3-100 chars, unique)
  - Store slug (3-50 alphanumeric)
  - Vendor type (individual/business)
  - Commission model (margin/markup)
  - Commission rate (0-1)
  
  Steps:
  1. Validate input
  2. Create vendor record (status: pending)
  3. Generate unique slug
  4. Create vendor wallet
  5. Send verification email
  
✓ UploadDocument(vendorId, docType, file)
  Supported types:
  - For Individual: NIC/Passport, Bank proof, Email/Phone OTP
  - For Business: BR cert, Form 1/20, TIN/VAT, Director NIC, Bank proof
  
  Validation:
  - PDF format only
  - Max size: 10MB
  - Virus scan
  - Store in S3
  
✓ SubmitDocuments(vendorId, docList)
  → Mark vendor as submitted
  → Notify admin for review
  → Trigger KYC verification workflow
```

**Document Management:**
```
✓ GetDocuments(vendorId)
  → Return all uploaded documents with status (pending/verified/rejected)
  
✓ UpdateDocumentStatus(docId, status, comment)
  → Admin only
  → Can reject with reason
  → Trigger re-upload notification
```

### 2.2 Vendor Profile Management

**Profile Operations:**
```
✓ GetVendorProfile(vendorId)
  → Return: store name, logo, description, status, commission config
  → With eager load: user, bank details, wallet, documents
  
✓ UpdateVendorProfile(vendorId, request)
  Vendor can update:
  - Store name (subject to admin approval)
  - Logo URL (S3)
  - Description
  - Contact details (phone, email)
  - Pickup address
  
  Cannot update:
  - Bank details (admin only after initial setup)
  - Commission rate (admin only)
  - Legal identity details (admin only)
  
✓ GetBankDetails(vendorId)
  → Admin only
  → Return encrypted account details
  
✓ UpdateBankDetails(vendorId, accountName, accountNumber, bankName, branch)
  → Admin only after initial setup
  → Log all changes for audit
  → Cannot be edited by vendor post-activation
  
✓ ValidateBankDetails(vendorId)
  → Required for vendor activation
  → Perform bank account validation
```

### 2.3 Vendor Agreement Management

**Agreement Flow:**
```
✓ PublishAgreement(version, fileUrl, title)
  → Admin only
  → Store new version
  → Mark previous version as inactive
  → Notify all vendors to re-accept
  
✓ GetAgreement(versionId)
  → Return agreement details with markdown/PDF
  
✓ AcceptAgreement(vendorId, agreementId)
  → Log acceptance with timestamp
  → Log IP address
  → Store accepted version reference
  → Unlock product listing if KYC approved
  
✓ CheckAgreementCompliance(vendorId)
  → Return true if vendor accepted latest agreement
  → Used to block product listing if outdated
  
✓ GetLatestAgreement()
  → Return current active agreement
```

### 2.4 Vendor Moderation & Lifecycle

**Status Management:**
```
Status Flow: Pending → Approved → Active
            → Suspended (from Active)
            → Rejected (from Pending)

✓ ApproveVendor(vendorId, approvalNote)
  → Admin only
  → Validates: Documents approved + Agreement accepted + Bank details set
  → Change status: pending → approved
  → Create default vendor wallet
  → Send confirmation email
  → Vendor can now list products
  
✓ RejectVendor(vendorId, rejectionReason)
  → Admin only
  → Change status: pending → rejected
  → Send rejection email with reason
  → Vendor can reapply
  
✓ SuspendVendor(vendorId, suspensionReason)
  → Admin only
  → Change status: approved → suspended
  → Automatically pause all products (status: draft)
  → Pause order processing
  → Send notification to vendor
  
✓ ReinstateVendor(vendorId)
  → Admin only
  → Change status: suspended → approved
  → Reactivate all products
  → Send reinstatement notification
  
✓ ListVendors(status, limit, offset)
  → Admin endpoint
  → Filter by status (pending/approved/suspended/rejected)
  → Return with user info, document count, wallet balance
  
✓ SearchVendors(query, filters)
  → Admin endpoint
  → Search by: store name, email, phone, vendor type
  → Filter by: status, commission model, date range
```

### 2.5 Vendor Staff Management

**Staff Account Operations:**
```
✓ CreateStaffAccount(vendorId, request)
  Fields: email, firstName, lastName, role (store_manager/staff)
  
  Generate temporary password
  Send invitation email
  Staff must change password on first login
  
✓ ListStaff(vendorId)
  → Return all staff under vendor
  → Show: email, role, status (active/inactive), created_at
  
✓ UpdateStaffRole(staffId, newRole)
  → Vendor can manage own staff
  → Available roles: store_manager (full access), staff (limited)
  
✓ DisableStaff(staffId)
  → Revoke access
  → Staff user marked inactive
  → All sessions invalidated
  
✓ EnableStaff(staffId)
  → Reactivate staff account
```

---

## 3. PRODUCT CATALOG & INVENTORY

### 3.1 Category & Brand Management

**Category Operations:**
```
✓ CreateCategory(name, slug, parentId, minPrice, maxPrice, maxDiscountPct, commissionRate, commissionModel)
  → Admin only
  → Support hierarchical: category → subcategory → sub-subcategory
  → Max 3 levels deep
  → Slug must be unique
  
✓ GetCategory(categoryId)
  → Return: name, slug, pricing rules, commission config, status
  
✓ ListCategories(parentId)
  → Return all top-level categories
  → If parentId provided, return subcategories
  → Include pricing rules and commission info
  
✓ UpdateCategory(categoryId, request)
  → Admin only
  → Can update: name, pricing rules, commission settings
  
✓ DeleteCategory(categoryId)
  → Admin only
  → Mark as inactive (soft delete)
  → Check no products in category before deletion
  
✓ SetCategoryCommission(categoryId, commissionModel, commissionRate)
  → Admin only
  → Override for entire category
  → Applies to all products in category
  → Highest priority in commission hierarchy
```

**Brand Operations:**
```
✓ CreateBrand(name, slug, logoUrl)
  → Admin only
  → Store logo in S3
  
✓ ListBrands()
  → Return all brands with logo URLs
  
✓ UpdateBrand(brandId, name, logoUrl)
  → Admin only
  
✓ GetBrandProducts(brandId, limit, offset)
  → Return paginated products by brand
```

### 3.2 Product Management (Vendor)

**Product CRUD:**
```
✓ CreateProduct(vendorId, request)
  Fields:
  - name (3-255 chars)
  - SKU (3-50 chars, unique per vendor)
  - description (10-5000 chars)
  - categoryId (required)
  - brandId (optional)
  - basePrice (> 0)
  - costPrice (>= 0, < basePrice)
  - weight (> 0)
  - dimensions (JSON: length, width, height)
  - warranty (string)
  - seoTitle, seoDescription
  - isReturnable (bool)
  - returnWindowDays (if returnable)
  
  Validation:
  - Vendor must be approved + agreed to latest agreement
  - Check pricing against category min/max
  - SKU uniqueness
  - Cost < base price
  
  Status: draft → ready for variants
  
✓ GetProduct(productId)
  → Return all product details + images + variants + reviews
  → Eager load: vendor, category, brand
  
✓ ListVendorProducts(vendorId, limit, offset)
  → Vendor endpoint
  → Filter by: status (draft/pending/published/archived)
  → Sort by: created_at, name, price
  
✓ UpdateProduct(productId, request)
  Fields: same as create (optional)
  
  Business rules:
  - If status = published → auto reset to pending_approval
  - If draft → can update freely
  - Vendor can only update own products
  
✓ ArchiveProduct(productId)
  → Soft delete
  → Status: published → archived
  → Still visible in order history
  → Not shown in listings
  
✓ SubmitProductForApproval(productId)
  → Change status: draft → pending_approval
  → Notify admin
```

### 3.3 Product Variants & Inventory

**Variant Management:**
```
✓ CreateVariant(productId, request)
  Fields:
  - SKU (unique)
  - price (> 0)
  - stock (>= 0)
  - attributes (JSON: {color: "red", size: "M"})
  
  Validation:
  - SKU must be unique across all variants
  - Price > 0
  
✓ GetVariants(productId)
  → Return all variants with stock levels
  → Sort by: attributes, price
  
✓ UpdateVariant(variantId, price, stock, attributes)
  → Vendor can update own variants
  → If price changed → check against category pricing rules
  
✓ UpdateVariantStock(variantId, newStock)
  → Direct stock update (manual)
  → Used by inventory sync or manual adjustment
  
✓ DecrementStock(variantId, quantity)
  → Decrease by quantity (for order placement)
  → Used in checkout
  → Check stock before decrement
  
✓ IncrementStock(variantId, quantity)
  → Increase by quantity (for refunds/returns)
  
✓ GetLowStockProducts(vendorId, threshold)
  → Return products with stock < threshold
  → Vendor dashboard alert
```

### 3.4 Product Images & SEO

**Image Management:**
```
✓ UploadProductImage(productId, variantId, file)
  → Store in S3
  → Generate thumbnail
  → Return image URL
  
✓ GetProductImages(productId)
  → Paginated list with sort_order
  → Include primary image indicator
  
✓ SetPrimaryImage(imageId)
  → Set which image shows in listing
  → Only one primary per product
  
✓ DeleteImage(imageId)
  → Remove from S3
  → Remove from database
  
✓ ReorderImages(productId, imageIds[])
  → Update sort_order for all images
```

**SEO Management:**
```
✓ SetSEOMetadata(productId, title, description, keywords)
  → Vendor can set
  → Used for meta tags
  → Max 160 chars for title, 160 for description
```

### 3.5 Product Approval (Admin)

**Approval Workflow:**
```
Status: draft → pending_approval → published (approved)
                                → archived (rejected)

✓ ListPendingProducts(limit, offset)
  → Admin endpoint
  → Filter by: category, vendor, date range
  → Show: name, vendor, category, created_at, pending duration
  
✓ ApproveProduct(productId, approvalNote)
  → Admin only
  → Check: all images present, pricing valid, description complete
  → Change status: pending_approval → published
  → Send notification to vendor
  
✓ RejectProduct(productId, rejectionReason)
  → Admin only
  → Change status: pending_approval → rejected
  → Send rejection email with reason
  → Vendor must fix and resubmit
  
✓ GetProductApprovalHistory(productId)
  → Return all approval/rejection records
  → Show: admin name, action, reason, timestamp
```

### 3.6 Product Bulk Operations

**Bulk Import:**
```
✓ BulkImportProducts(vendorId, csvFile)
  CSV columns: name, SKU, description, category, price, cost_price, stock, variantAttributes
  
  Process:
  1. Parse CSV
  2. Validate each row
  3. Check SKU uniqueness
  4. Check pricing rules
  5. Create all products (batch insert)
  6. Return: success count, error count with details
  
✓ GetBulkImportStatus(importId)
  → Return import progress and any errors
```

**Bulk Export:**
```
✓ ExportProducts(vendorId)
  → Generate CSV with all vendor products
  → Include: SKU, name, price, cost_price, stock, variants, status
  → Send download link
```

---

## 4. PRICING, PROMOTIONS & DISCOUNTS

### 4.1 Pricing Governance

**Price Rules Management:**
```
✓ SetCategoryPricingRules(categoryId, minPrice, maxPrice, maxDiscountPct)
  → Admin only
  → Enforced on all products in category
  → Block checkout if violated
  
✓ GetCategoryPricingRules(categoryId)
  → Return min/max prices and max discount allowed
  
✓ ValidateProductPrice(productId, basePrice, proposedDiscount)
  → Check against category rules
  → Return: isValid, violatedRules[]
  
✓ EnforcePrice ViolationBlock(productId)
  → Block product listing if price violated
  → Send admin alert
```

### 4.2 Promotion Management

**Admin Promotions (Platform-funded):**
```
✓ CreatePromotion(request)
  Fields:
  - type (product_discount, coupon, bogo, bundle, order_discount, free_shipping, bank_ipg)
  - code (for coupon, max 50 chars)
  - discountType (flat/percent)
  - discountValue (amount or percentage)
  - minOrderValue (optional)
  - validFrom, validTo (datetime)
  - usageLimit (optional)
  - applicableCategories (list of category IDs)
  - fundingType (vendor/platform/shared)
  - vendorSharePct (if shared, percentage for vendor)
  
✓ GetPromotion(promotionId)
  → Return promotion details with usage stats
  
✓ ListActivePromotions(categoryId)
  → Return promotions applicable to category
  → Filter by: valid date range, usage limits
  
✓ UpdatePromotion(promotionId, request)
  → Can update: validity period, usage limit, discount value
  → Cannot change: funding type, code
  
✓ DeactivatePromotion(promotionId)
  → Set status: inactive
  → No new orders can use it
  
✓ GetPromotionStats(promotionId)
  → Return: usage count, total discount given, affected revenue
```

**Vendor Promotions (Vendor-funded):**
```
✓ CreateVendorPromotion(vendorId, request)
  Vendor can create:
  - Product-level discounts
  - Coupon codes
  - BOGO offers
  - Bundle offers
  
  Constraints:
  - Cannot exceed category maxDiscountPct
  - Must not violate pricing rules
  
✓ ListVendorPromotions(vendorId)
  → Return vendor's active promotions
  
✓ UpdateVendorPromotion(promotionId, request)
  → Vendor can manage own
```

### 4.3 Discount Application Logic

**Promotion Priority & Calculation:**
```
Priority order (highest to lowest):
1. Bank IPG promotions
2. Product-level / bundle offers
3. Coupon codes

✓ ApplyPromotions(orderId, promoCodes[])
  → For each promotion:
    1. Check validity (date, usage limit, min order value)
    2. Get applicable items
    3. Calculate discount amount
    4. Update order.discountTotal
  
  → Calculate based on:
    - Discount type: flat amount or percentage
    - Discount funding: vendor vs platform
    - Sequential application (each applies to reduced amount)
  
  → Return: totalDiscount, breakdownByPromotion, finalPrice
  
✓ ValidatePromoCode(code)
  → Check: exists, active, usage limit not exceeded, date valid
  → Return: isValid, promotion details, applicableItems
  
✓ IncrementPromotionUsage(promotionId)
  → Increment used_count
  → Check if usage_limit exceeded
  
✓ CalculateDiscountAmount(baseAmount, promotion)
  → If discountType = flat: return promotion.discountValue
  → If discountType = percent: return baseAmount * (promotion.discountValue / 100)
  → Cap by minOrderValue constraint
```

---

## 5. CART & CHECKOUT

### 5.1 Shopping Cart

**Cart Operations:**
```
✓ CreateCart(userId or sessionId)
  → For logged-in users: link to user
  → For guests: create session-based cart
  → Set created_at timestamp
  
✓ GetCart(userId or sessionId)
  → Return: items, subtotal, line items
  → Eager load: product, variant, vendor info
  → Recalculate prices (in case changed)
  
✓ AddToCart(cartId, variantId, quantity)
  Fields: productVariantId, quantity (1-1000)
  
  Validation:
  - Variant exists
  - quantity > 0 and <= 1000
  - Stock available (qty <= variant.stock)
  - Product status = published
  - Vendor status = approved
  
  Logic:
  - If same variant already in cart → increment qty
  - Set unit_price = current variant.price
  - cart.updated_at = now
  
✓ UpdateCartItem(cartItemId, newQuantity)
  Validation: same as AddToCart
  
  Logic:
  - Update qty
  - Check still in stock
  - Recalculate subtotal
  
✓ RemoveFromCart(cartItemId)
  → Delete cart item
  → Update cart.updated_at
  
✓ ClearCart(cartId)
  → Delete all items in cart
  → Keep cart record (for abandon tracking)
  
✓ GetCartSummary(cartId)
  → Return:
    - itemCount
    - subtotal
    - distinctVendors
    - estimatedShipping (placeholder)
    - applicablePromotions
  
✓ RecalculateCartPrices(cartId)
  → Called before checkout
  → Check each item still exists and published
  → Update unit prices if changed
  → Check stock availability
  → Return: updated prices, any items removed due to unavailable
```

### 5.2 Checkout Process

**Checkout Flow:**
```
✓ ValidateCheckout(cartId, request)
  Request fields:
  - shippingAddressId or shippingAddress (new)
  - paymentMethod (bank_ipg/emi/cod/manual_transfer)
  - promoCodes[]
  
  Validations:
  1. Cart not empty
  2. All items still in stock
  3. All vendors still approved
  4. Product prices haven't violated rules
  5. Shipping address valid (for delivery)
  6. Promo codes valid
  
  Return: isValid, violations[], estimatedTotal
  
✓ Checkout(cartId, request)
  CRITICAL FUNCTION - Main order creation logic
  
  Steps:
  1. ValidateCheckout
  2. CreateMasterOrder
     - status: pending
     - paymentMethod: from request
     - paymentStatus: pending
     - shippingAddress: from request
  
  3. ApplyPromotions (get discount breakdown)
  
  4. **SplitOrder()** - CRITICAL
     - Group cart items by vendor
     - For each vendor group:
       a. Calculate subtotal
       b. Calculate commission per item (CommissionService)
       c. Create SubOrder:
          - vendor_id
          - status: pending
          - subtotal
          - commission_amount
          - vendor_earning = subtotal - commission
       d. Create OrderItems linked to SubOrder
     
  5. UpdateOrderTotals
     - order.subtotal = sum all suborders.subtotal
     - order.discountTotal = from promotions
     - order.shippingTotal = calculated by shipping service
     - order.grandTotal = subtotal - discount + shipping
  
  6. DecrementStock
     - For each item, decrement variant.stock
     - If decrement fails → rollback entire order
  
  7. CreateWalletTransactions (for each vendor)
     - transaction type: order
     - amount: vendor_earning
     - status: pending (until order completes)
  
  8. SendNotifications
     - Vendor: new order email with details
     - Customer: order confirmation
  
  9. ClearCart
  
  10. Return: order, suborders, confirmation details
  
  Error handling:
  - Stock depletion → return error, don't create order
  - Promo invalid → return error
  - Vendor suspended → return error
  - Payment method validation fails → return error
```

**Promo Code Application:**
```
✓ ApplyCoupon(cartId, promoCode)
  → Validate promo code
  → Calculate discount
  → Update cart.appliedPromoCodes[]
  → Recalculate cart total
  
✓ RemoveCoupon(cartId, promoCode)
  → Remove from cart.appliedPromoCodes
  → Recalculate total
```

### 5.3 Payment Processing

**Payment Methods (MVP Support):**
```
✗ ProcessBankIPG(orderId, cardDetails)
  → Integrate with bank gateway
  → Redirect to payment page
  → Handle callback (webhook)
  → Update order.paymentStatus
  
✗ ProcessEMI(orderId, emiPlan)
  → Check EMI eligibility
  → Redirect to EMI provider
  
✗ ProcessCOD(orderId)
  → Set paymentStatus: pending
  → Wait for delivery personnel to collect
  → Mark as paid after delivery
  
✓ ProcessManualTransfer(orderId, transferProof)
  → For MVP: accept screenshot/reference
  → Admin manual verification
  → Update order.paymentStatus
  
✓ GetPaymentStatus(orderId)
  → Return current payment_status
  → If pending > 24h → send reminder email
  
✓ MarkPaymentPaid(orderId)
  → Admin marks as paid (for manual verification)
  → Update order.paymentStatus: paid
  → Trigger order processing
```

---

## 6. ORDER PROCESSING & FULFILLMENT

### 6.1 Order Management

**Order Operations:**
```
✓ GetOrder(orderId, userId)
  → Return complete order with all sub-orders and items
  → Validate user ownership (customer/vendor/admin)
  → Eager load: suborders, items, vendor info, shipping
  
✓ ListCustomerOrders(customerId, limit, offset)
  → Customer endpoint
  → Filter by: status, date range
  → Sort by: created_at desc
  → Return with status summary
  
✓ ListVendorOrders(vendorId, limit, offset)
  → Vendor endpoint
  → Return only their sub-orders
  → Filter by: status, date range
  → Include: customer info, items, totals
  
✓ ListAllOrders(filters)
  → Admin endpoint
  → Filter by: status, vendor, customer, date range, payment status
  → Return with: GMV calculation, commission details
  
✓ SearchOrders(query)
  → Admin endpoint
  → Search by: order ID, customer email, vendor name, SKU
  
✓ GetOrderTimeline(orderId)
  → Return all status changes with timestamp
  → Show: who changed, from/to status, reason/note
```

### 6.2 Order Status Management

**Status Transitions:**
```
Master Order: pending → paid → processing → shipped → delivered → cancelled/refunded

SubOrder: pending → ready_to_ship → shipped → delivered → cancelled

✓ UpdateOrderStatus(orderId, newStatus, note)
  → Admin only
  → Can transition to any status
  → Log all changes
  → Send notifications
  
✓ UpdateSubOrderStatus(subOrderId, newStatus, note)
  → Vendor: pending → ready_to_ship ONLY
  → Admin: any transition
  → Log changes
  → Notify customer via email/SMS
  
✓ CancelOrder(orderId, reason)
  → Customer can cancel if pending/processing
  → Admin can cancel any status
  → Actions:
    1. If payment_status = paid → process refund
    2. Increment stock for all items
    3. Reverse wallet transactions
    4. Cancel all sub-orders
    5. Send notification to vendor
  
✓ CancelSubOrder(subOrderId, reason)
  → Vendor can cancel if pending
  → Admin can cancel any status
  → Actions:
    1. Revert stock for items in this sub-order
    2. Adjust vendor earnings
    3. Adjust order total
  
✓ MarkAsShipped(subOrderId, trackingNumber)
  → Vendor updates
  → Set tracking number
  → Change status: ready_to_ship → shipped
  → Notify customer with tracking info
  
✓ MarkAsDelivered(subOrderId)
  → Admin or delivery personnel
  → Check all sub-orders shipped
  → Update master order: shipped → delivered
  → Mark payment as received (for COD)
  → Trigger commission settlement
  
✓ SetShippingAddress(orderId, address)
  → Customer can update if still pending
  → Cannot update after shipped
```

### 6.3 Order Status Logging

**Audit Trail:**
```
✓ CreateStatusLog(orderId, subOrderId, fromStatus, toStatus, changedBy, note)
  → Log every status change
  → Store in order_status_logs table
  → Include: timestamp, user role, reason
  
✓ GetStatusHistory(orderId)
  → Return chronological status changes
  → Show: status, timestamp, who changed, reason
```

### 6.4 Shipping & Delivery

**Shipping Calculation:**
```
✗ CalculateShippingCost(shippingAddress, items)
  → Integration with logistics partners (GHN, GHTK)
  → Based on: destination, weight, zone
  → Return: shipping cost, estimated days, partner
  
✗ GetShippingZones()
  → Return all defined shipping zones with rates
  → Admin configurable
  
✗ CreateShipment(subOrderId)
  → Call logistics API
  → Create pickup/delivery orders
  → Store tracking number
  
✗ GetShipmentStatus(trackingNumber)
  → Call logistics API
  → Return current status and location
  → Cache in Redis (5 min TTL)
  
✗ RequestPickup(subOrderId, pickupAddress)
  → Vendor requests pickup
  → Schedule pickup with logistics partner
  → Send pickup confirmation to vendor
  
✗ TrackShipment(orderId)
  → Customer endpoint
  → Return shipment status and ETA
  → For all sub-orders
```

---

## 7. RETURNS & REFUNDS

### 7.1 Return Request Management

**Return Workflow:**
```
Status: requested → approved/rejected → picked_up → completed

✓ InitiateReturn(orderId, subOrderId, reason, evidenceUrls[])
  Fields:
  - orderId: master order
  - subOrderId: which sub-order (optional, if all items)
  - reason: predefined list or custom
  - evidenceUrls: customer photos/videos (max 5)
  
  Validation:
  - Order delivered (status = delivered)
  - Within return window (now - delivery_date <= returnWindowDays)
  - Product marked as returnable
  - Not already returned
  
  Actions:
  1. Create return_request record
  2. Set status: requested
  3. Notify vendor for review
  4. Notify customer: request received
  
✓ ListReturns(customerId or vendorId)
  → Customer: show all return requests for their orders
  → Vendor: show requests for their items
  → Filter by: status, date range
  
✓ GetReturnRequest(returnId, userId)
  → Return full details with evidence
  → Validate user access
  
✓ ApproveReturn(returnId, approvalNote)
  → Vendor or Admin
  → Set status: requested → approved
  → Log who approved
  → Send approval email to customer with pickup instructions
  
✓ RejectReturn(returnId, rejectionReason)
  → Vendor or Admin
  → Set status: requested → rejected
  → Send rejection email with reason
  → Cannot refund rejected returns
  
✓ MarkAsPickedUp(returnId)
  → Admin or delivery person
  → Set status: approved → picked_up
  → Verify item received
  
✓ CompleteReturn(returnId)
  → Admin marks return complete
  → Set status: picked_up → completed
  → Triggers refund processing
```

### 7.2 Refund Processing

**Refund Flow:**
```
✓ ProcessRefund(returnId)
  Called after return marked complete
  
  Logic:
  1. Get return request with amount
  2. Determine refund method (original payment / bank transfer)
  3. If original method available:
     - Call payment gateway reverse
  4. Else:
     - Create bank transfer request
     - Admin manually processes
  5. Create refund record
  6. Set status: pending → processed
  7. Update order: refund_total
  8. Update wallet: reverse vendor earnings
  9. Reverse loyalty points earned
  10. Send refund notification
  
✓ GetRefundStatus(refundId)
  → Return current status (pending/processed/failed)
  → Show: amount, method, processed_at
  
✓ ListRefunds(vendorId or customerId)
  → Vendor: refunds for their items
  → Customer: refunds for their orders
  → Show: amount, status, date
```

### 7.3 Financial Impact of Returns

**Wallet & Commission Reversal:**
```
✓ ReverseVendorEarnings(orderId)
  → Called on refund approval
  → Get vendor_earning from order
  → Create reverse wallet transaction
  → Deduct from vendor.balance
  
✓ ReverseLoyaltyPoints(customerId, points)
  → Called on refund
  → Deduct points from customer
  
✓ ApplyRefundRules(order, refundAmount)
  Rules:
  - Platform-funded discount: not reclaimed from vendor
  - Vendor-funded discount: deducted from vendor earnings
  - Shared discount: split reversal between vendor and platform
```

---

## 8. COMMISSION, SETTLEMENT & FINANCIAL LOGIC

### 8.1 Commission Calculation

**Commission Service (Core Business Logic):**
```
✓ CalculateCommission(orderItem, vendor, category)
  Priority hierarchy:
  1. Category commission (highest)
  2. Vendor commission
  3. Platform default (5% margin)
  
  Models:
  
  **Margin Model:**
  - commission = (selling_price - cost_price) * rate
  - Example: $100 price, $70 cost, 20% rate
    → commission = ($100 - $70) * 0.20 = $6
  
  **Markup Model:**
  - commission = selling_price / (1 + rate) * (1 - 1/(1 + rate))
  - Example: $100 price, 10% markup
    → commission ≈ $9.09
  
  Return: {
    commissionAmount,
    commissionModel,
    commissionRate,
    vendorEarning = subtotal - commission
  }
  
✓ ValidateCommissionRate(rate, model)
  → Ensure 0 <= rate <= 1
  → Ensure margin rate >= 0
  → Validate against category max rate
  
✓ GetApplicableCommission(productId, vendorId)
  → Check category commission first
  → Fall back to vendor commission
  → Fall back to platform default
  → Return applicable rate and model
  
✓ LockCommission(orderId)
  → After order confirmation
  → Save commission_model_snapshot in order_items
  → Prevent future changes from affecting this order
```

### 8.2 Discount Impact on Commission

**Discount Funding Models:**
```
✓ CalculateCommissionWithDiscount(item, discount)
  
  If discount.fundingType = vendor:
    → Commission calculated on (price - discount)
    → Vendor pays for discount
    → Example: $100 price, 10% vendor discount
      - Discounted price: $90
      - Commission on $90
      - Vendor earning reduced by $10
  
  If discount.fundingType = platform:
    → Commission calculated on original price
    → Platform pays for discount
    → Vendor unaffected
  
  If discount.fundingType = shared:
    → Split discount between vendor and platform
    → Commission calculated based on shared amount
    → Share ratio: discount.vendorSharePct
```

### 8.3 Vendor Wallet & Transactions

**Wallet Management:**
```
✓ CreateVendorWallet(vendorId)
  → Auto-created on vendor approval
  → Fields: balance, pending_balance
  → Initial: 0
  
✓ GetVendorWallet(vendorId)
  → Return: balance, pending_balance, lastUpdated
  
✓ CreateWalletTransaction(vendorId, transactionData)
  Fields:
  - orderId (if order-related)
  - refType (order/return/settlement/adjustment)
  - amount (can be negative)
  - type (credit/debit)
  - description
  
  Logic:
  - Create transaction record (immutable)
  - Calculate balance_after = balance + amount (if credit) or - amount (if debit)
  - Log for audit
  
✓ GetWalletLedger(vendorId, limit, offset)
  → Vendor endpoint
  → Return paginated transactions
  → Filter by: date range, type (credit/debit), ref_type
  → Show: amount, type, balance_after, description, date
  
✓ GetWalletBalance(vendorId)
  → Calculate current balance from all transactions
  → Cache in Redis (5 min TTL)
  
✓ DeductFromWallet(vendorId, amount, reason)
  → Admin only (for adjustments)
  → Create debit transaction
  → Update balance
  → Log reason for audit
```

### 8.4 Settlement Process

**Vendor Payout Settlement:**
```
✓ InitiateSettlement(vendorId, periodStart, periodEnd)
  → Admin triggers
  → Calculate:
    1. Gross sales: sum of all order subtotals in period
    2. Commission deducted: sum of commission_amount
    3. Refund deductions: sum of refunded amounts
    4. Net payable: gross - commission - refunds
  
  → Create settlement record
  → Set status: pending
  → Send to vendor for review (notification)
  
✓ ListSettlements(vendorId, limit, offset)
  → Vendor endpoint
  → Show: period, gross_sales, commission, refunds, net_payable, status
  → Filter by: status, date range
  
✓ GetSettlementDetails(settlementId)
  → Return breakdown:
    - Total orders count
    - Total gross sales
    - Commission deducted (itemized)
    - Refunds deducted
    - Final net payable
    - Recommended payout bank account
  
✓ ApproveSettlement(settlementId, adminNote)
  → Admin reviews and approves
  → Set status: pending → processing
  → Send to payment processing system
  
✓ MarkSettlementAsPaid(settlementId, paymentReference)
  → After actual bank transfer
  → Set status: processing → paid
  → Update paid_at timestamp
  → Create wallet transaction (debit for settlement)
  → Send confirmation to vendor
  
✓ GetSettlementCalculation(vendorId, periodStart, periodEnd)
  → Preview settlement (before initiating)
  → Show what will be settled
  → Useful for vendor estimation
  
✓ GetSettlementSchedule()
  → Return configured settlement cycle
  → Example: weekly (every Monday), bi-weekly, monthly
  → Admin configurable
```

### 8.5 Financial Records & Audit

**Audit Logging:**
```
✓ CreateAuditLog(userId, entityType, entityId, action, beforeData, afterData, ip)
  → Log all financial operations
  → Store before/after snapshots for audit trail
  → Include: admin name, action, timestamp, IP
  
✓ GetAuditLogs(filters)
  → Admin endpoint
  → Filter by: entity type, action, date range, user
  → Return: immutable log entries
  
✓ GetFinancialReport(vendorId, periodStart, periodEnd)
  → Comprehensive financial statement
  → Include: sales, commissions, refunds, settlements
  → Export as PDF/Excel
  
✓ ReconcileFinancials(vendorId)
  → Admin tool to verify accuracy
  → Sum order amounts vs settlement amounts
  → Detect discrepancies
  → Flag for review if variance > threshold
```

---

## 9. ADMINISTRATION & REPORTING

### 9.1 Admin Dashboard Metrics

**Financial Dashboard:**
```
✓ GetDashboardMetrics()
  Return:
  - Total GMV (Gross Merchandise Value)
  - Total commission earned (platform)
  - Pending payouts (to vendors)
  - Active vendors count
  - Active products count
  - Pending approvals count
  - Today's orders count
  
✓ GetGMVTrend(period: week/month/quarter)
  → Return GMV over time (for charts)
  → Compare with previous period
  
✓ GetCommissionBreakdown(period)
  → Breakdown by: vendor, category, commission model
  → Show: amount, percentage, trend
```

### 9.2 Reports

**Order Reports:**
```
✓ GenerateOrdersReport(filters)
  Filters: dateRange, vendor, customer, status, paymentStatus
  
  Return:
  - Total orders
  - Breakdown by status
  - Breakdown by vendor
  - GMV by vendor/category
  - Payment method distribution
  - Export as CSV/Excel
  
✓ GenerateSalesReport(dateRange)
  → GMV per day/week/month
  → Top performing vendors
  → Top selling products
  → Seasonal trends
  
✓ GenerateCommissionReport(dateRange)
  → Commission per vendor
  → Commission breakdown by category
  → Average commission rate
  → Platform earnings
  
✓ GenerateVendorPerformanceReport(vendorId)
  → Order count
  → Revenue
  → Commission paid
  → Return rate
  → Customer satisfaction (avg rating)
  → Compliance status
  
✓ GenerateInventoryReport()
  → Total stock by category
  → Low stock alerts
  → Stock turnover
  → Dead stock (no sales in 30 days)
  
✓ GenerateReturnsReport(dateRange)
  → Return count by vendor
  → Return rate (%)
  → Most returned products
  → Return reasons breakdown
  → Refund amounts
  
✓ GenerateSettlementReport(dateRange)
  → Settlements initiated
  → Settlements paid
  → Average settlement amount
  → Settlement by vendor
  → Outstanding payments
```

### 9.3 Role-Based Access Control (Admin Panel)

**Super Admin Functions:**
```
✓ CreateRole(name, permissions[])
✓ UpdateRolePermissions(roleId, permissions[])
✓ DeleteRole(roleId)
✓ AssignRoleToUser(userId, roleId)
✓ ListUsers(roleFilter)
✓ CreateAdminUser(email, firstName, lastName, roleId)
✓ DeactivateUser(userId)
```

**Admin Functions:**
```
✓ ApproveVendor(vendorId)
✓ RejectVendor(vendorId)
✓ SuspendVendor(vendorId)
✓ ApproveProduct(productId)
✓ RejectProduct(productId)
✓ ApproveReturn(returnId)
✓ RejectReturn(returnId)
✓ InitiateSettlement(vendorId)
✓ ViewAuditLogs()
✓ ViewReports()
```

---

## 10. CUSTOMER ENGAGEMENT & NOTIFICATIONS

### 10.1 Wishlist Management

**Wishlist Operations:**
```
✓ AddToWishlist(customerId, productId)
  → For logged-in users only
  → Check if already in wishlist (prevent duplicates)
  → Store in wishlist table with timestamp
  
✓ GetWishlist(customerId)
  → Return paginated wishlist items
  → Show: product details, current price, availability, vendor
  → Sort by: added_at desc
  
✓ RemoveFromWishlist(customerId, productId)
  → Delete wishlist entry
  
✓ IsInWishlist(customerId, productId)
  → Return boolean (for UI)
  
✓ GetWishlistCount(customerId)
  → Return count of items in wishlist
  
✓ NotifyPriceDropForWishlist(productId, newPrice)
  → Find all customers with this product in wishlist
  → Send notification: price dropped
  → Include: original price, new price, discount %
  → Link to product page
```

### 10.2 Star Ratings & Reviews

**Review System:**
```
✓ CreateReview(orderId, orderItemId, rating, reviewText)
  Constraints:
  - Only verified purchasers (order.status = delivered)
  - One review per item
  - Rating: 1-5 stars
  - Review text: optional, 10-5000 chars
  - Can include photos (max 5)
  
  Process:
  1. Create review record
  2. Set status: pending_approval
  3. Notify admin for moderation
  
✓ GetProductReviews(productId, sortBy)
  → Public endpoint
  → Filter: approved reviews only
  → Sort: helpful_count desc, rating asc, recent
  → Show: rating, text, reviewer name (hidden), date, helpful_count
  
✓ GetAverageRating(productId)
  → Calculate: average rating (1-5)
  → Return: avgRating, totalReviews, ratingDistribution
  → Cache in Redis (1 hour TTL)
  
✓ ApproveReview(reviewId, approvalNote)
  → Admin only
  → Set status: pending_approval → published
  → Notify reviewer: review published
  
✓ RejectReview(reviewId, rejectionReason)
  → Admin only
  → Set status: pending_approval → rejected
  → Notify reviewer: reason for rejection
  
✓ HideReview(reviewId, adminReason)
  → Admin can hide published reviews
  → Don't delete (audit trail)
  → Just set visible: false
  
✓ MarkHelpful(reviewId, userId)
  → Customer marks review as helpful
  → Increment helpful_count
  → Prevent duplicate marks (one per user)
  
✓ ReportReview(reviewId, reason)
  → Customer can report inappropriate reviews
  → Create report record
  → Notify admin for review
```

### 10.3 Loyalty Points

**Points System:**
```
✓ EarnPoints(orderId, customerId, amount)
  → Called after order completion (post-return window)
  → Rules configurable by admin
  → Example: 1 point per LKR 100 spent
  
  Calculation:
  - Base points: order.grandTotal / configuredDenomination
  - Bonus multiplier: promotions, loyalty tiers
  - Final points: basePoints * multiplier
  
  Process:
  1. Create points transaction (credit)
  2. Update customer.totalPoints
  3. Send email: points earned
  
✓ RedeemPoints(customerId, pointsToRedeem, orderId)
  → Called at checkout
  → Check sufficient balance
  → Constraints:
    - Max redeemable per order: configurable
    - Cannot use on restricted categories
    - Cannot use with conflicting promotions
  
  Process:
  1. Create points transaction (debit)
  2. Calculate discount: points * redeemRate (e.g., 1 point = LKR 1)
  3. Apply as order discount
  4. Update customer.totalPoints
  
✓ ExpirePoints(customerId)
  → Scheduled job (monthly)
  → Expire points older than configurable period (e.g., 1 year)
  → Create debit transaction with reason: "expired"
  
✓ GetPointsBalance(customerId)
  → Return: total points, points expiring soon
  
✓ GetPointsHistory(customerId, limit, offset)
  → Return paginated transaction log
  → Show: type (earned/redeemed/expired), amount, date, reason, balance_after
  
✓ AdjustPoints(customerId, amount, reason)
  → Admin only
  → Manual adjustment for customer service
  → Log reason for audit
  
✓ GetPointsExchangeRate()
  → Return current exchange rate
  → Admin configurable
  → Example: 100 points = LKR 100
```

### 10.4 Email Notifications

**Customer Notifications:**
```
✓ SendOrderConfirmation(orderId)
  → Sent after order creation
  → Include: order details, items, tracking link, estimated delivery
  
✓ SendPaymentConfirmation(orderId)
  → Sent after payment processed
  
✓ SendShipmentNotification(subOrderId, trackingNumber)
  → Sent when item shipped
  → Include: tracking number, carrier, estimated delivery date
  
✓ SendDeliveryNotification(orderId)
  → Sent when delivered
  → Include: delivery date, return window expiry
  
✓ SendReturnApprovalNotification(returnId)
  → Sent when return approved
  → Include: pickup instructions, return address
  
✓ SendRefundNotification(refundId)
  → Sent when refund processed
  → Include: refund amount, method, expected time to receive
  
✓ SendPointsEarnedNotification(customerId, pointsEarned)
  → Sent after order completes
  → Include: points earned, total balance, redemption rate
  
✓ SendAbandonedCartReminder(customerId)
  → Scheduled: 24 hours after cart abandoned
  → Include: cart items, prices, discount reminder
```

**Vendor Notifications:**
```
✓ SendNewOrderNotification(vendorId, orderId)
  → Sent immediately on order creation
  → Include: order details, customer info, items, deadline to ship
  
✓ SendOrderStatusChangeNotification(vendorId, subOrderId, newStatus)
  → Sent on any admin-initiated status change
  → Include: from/to status, reason, action required
  
✓ SendReturnRequestNotification(vendorId, returnId)
  → Sent when customer initiates return for their item
  → Include: product details, reason, evidence links, action required
  
✓ SendSettlementNotification(vendorId, settlementId)
  → Sent when settlement initiated
  → Include: period, gross sales, commission, net payable, deadline for objection
  
✓ SendSettlementPaidNotification(vendorId, settlementId)
  → Sent when settlement paid
  → Include: amount, payment reference, bank transfer details
  
✓ SendDocumentApprovalNotification(vendorId, docId)
  → Sent when KYC document approved/rejected
  → If rejected: include reason and reupload link
  
✓ SendProductApprovalNotification(vendorId, productId, approved)
  → Sent when product approved/rejected
  → If rejected: include reason and edit link
  
✓ SendVendorSuspensionNotification(vendorId, reason)
  → Sent when vendor suspended
  → Include: reason, action items to resolve
```

### 10.5 WhatsApp Integration

**WhatsApp Support Channel:**
```
✗ InitializeWhatsAppChat(orderId)
  → Provide WhatsApp contact button on order page
  → Link to support WhatsApp number
  → Pre-fill order ID in message
  
✗ SendWhatsAppTemplateMessage(recipientPhone, template, params)
  → Send predefined templates:
    - Order confirmation
    - Shipment tracking
    - Return approval
    - Refund status
  
✗ ReceiveWhatsAppMessage(messageData)
  → Webhook receiver
  → Parse incoming message
  → Route to support team
  → Log conversation
  
✗ SendWhatsAppReply(chatId, replyMessage)
  → Support team responds to customer inquiry
  → Message sent via WhatsApp API
```

---

## 11. ERP INTEGRATION

**Data Synchronization:**
```
✗ SyncProductsFromERP(erpConnectionId)
  → Scheduled job (configurable interval)
  → Pull product data from ERP API
  → Match by SKU
  → Update: name, price, inventory
  → Handle conflicts: ERP is source of truth
  → Log sync process
  
✗ SyncInventoryToERP(orderId)
  → After order placement
  → Push stock decrement to ERP
  → Handle sync failures with retry logic
  
✗ GetERPSyncStatus()
  → Return last sync time, status (success/failed)
  → Show any pending syncs
  
✗ ConfigureERPConnection(connectionDetails)
  → Admin configures ERP API endpoint
  → Store credentials securely
  → Test connection
  
✗ RetryFailedSync(syncId)
  → Manual retry for failed sync
  → Log retry attempt
```

---

## 12. SEARCH & DISCOVERY

**Product Search:**
```
✗ SearchProducts(query, filters)
  Public endpoint
  → Full-text search on: name, SKU, description
  → Filter by: category, brand, price range, rating, availability
  → Sort by: relevance, price, rating, newest
  → Return paginated results with facets
  
✗ GetProductComparison(productIds[])
  → Compare selected products
  → Show: specs, price, rating, reviews side-by-side
  → Limit: max 4 products
  
✗ GetRelatedProducts(productId)
  → Return similar products by: category, attributes, brand
  → Exclude out-of-stock
  → Limit: 10 results
  
✗ GetTrendingProducts(period, category)
  → Return best-selling products
  → Filter by: category, period (week/month)
  → Limit: 20 results
```

---

## SUMMARY: Function Count by Module

| Module | Functions Count | Status |
|--------|-----------------|--------|
| 1. Authentication | 8 | 80% |
| 2. Vendor Management | 25 | 50% |
| 3. Product Catalog | 35 | 40% |
| 4. Pricing & Promotions | 15 | 30% |
| 5. Cart & Checkout | 12 | 20% |
| 6. Order Processing | 20 | 30% |
| 7. Returns & Refunds | 10 | 10% |
| 8. Commission & Settlement | 20 | 60% |
| 9. Administration | 20 | 20% |
| 10. Customer Engagement | 25 | 0% |
| 11. ERP Integration | 5 | 0% |
| 12. Search & Discovery | 6 | 0% |
| **TOTAL** | **201 Functions** | **~30%** |

---

## CRITICAL PATH TO MVP (Priority)

**Must Have (Week 1-2):**
1. ✓ Auth (Register, Login, Token refresh)
2. ✓ Vendor Registration & KYC approval
3. ✓ Product CRUD & Approval workflow
4. ⚠️ **Checkout & Order Splitting** (SplitOrder)
5. ⚠️ **Commission Calculation**
6. ⚠️ **Payment Processing** (at least manual for MVP)
7. ⚠️ **Order Status Management**
8. ⚠️ **Vendor Notifications** (email)
9. ⚠️ **Refund Processing**
10. ⚠️ **Settlement** (vendor payouts)

**Should Have (Week 3-4):**
- Return request workflow
- Wishlist
- Basic search
- Admin reports
- Inventory sync

**Nice to Have (Phase 2):**
- WhatsApp integration
- Reviews & ratings
- Loyalty points
- ERP integration
- Advanced search (full-text)
