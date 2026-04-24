# Coolmate Multivendor eCommerce Backend

Go + Gin backend for a comprehensive multivendor eCommerce platform. Features RBAC, vendor onboarding with KYC, product management with approval workflow, order splitting per vendor, and financial settlement system.

## Quick Start

### Prerequisites
- Go 1.22+
- Docker & Docker Compose
- Git

### 1. Setup Environment

```bash
# Clone the repo
cd coolmate-prototype\ backend

# Create .env from example
cp .env.example .env

# Start PostgreSQL, Redis, and MinIO
docker-compose up -d
```

### 2. Run Migrations & Start Server

```bash
# Download dependencies
go mod tidy

# Run the server (auto-migrates models)
go run cmd/api/main.go
```

Server starts on `http://localhost:8080`

### 3. Test API

```bash
# Register customer
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "customer@example.com",
    "password": "password123",
    "firstName": "John",
    "lastName": "Doe",
    "phone": "1234567890"
  }'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "customer@example.com",
    "password": "password123"
  }'
```

---

## Architecture

### Directory Structure
```
├── cmd/api/main.go              # Entry point, route setup
├── internal/
│   ├── config/                  # Configuration management
│   ├── database/                # DB & Redis connections
│   ├── middleware/              # Auth, RBAC, CORS
│   ├── models/                  # GORM models
│   ├── handlers/                # HTTP request handlers
│   ├── services/                # Business logic layer
│   ├── repositories/            # Data access layer
│   └── utils/                   # Helpers (response, pagination, password)
├── pkg/
│   ├── auth/jwt.go              # JWT token generation/validation
│   ├── storage/s3.go            # S3/MinIO file upload
│   └── cache/redis_cache.go     # Redis caching
├── migrations/                  # SQL migration files (optional)
└── docker-compose.yml           # Local dev environment
```

### Layered Architecture

```
HTTP Request
    ↓
[Handler] - Parse request, call service
    ↓
[Service] - Business logic, orchestration
    ↓
[Repository] - Database queries via GORM
    ↓
PostgreSQL Database
```

- **Handlers**: Receive HTTP requests, validate input, call services, return responses
- **Services**: Contain business logic (commission calculations, order splitting, etc.)
- **Repositories**: Encapsulate database access with interfaces for testability
- **Models**: GORM struct definitions for all entities

---

## API Routes

### Authentication (Public)
```
POST   /api/v1/auth/register        - Register customer
POST   /api/v1/auth/login           - Login
POST   /api/v1/auth/refresh         - Refresh access token
POST   /api/v1/auth/logout          - Logout
```

### Vendor (Protected)
```
POST   /api/v1/vendor/register      - Vendor registration
POST   /api/v1/vendor/documents/upload - Upload KYC docs
GET    /api/v1/vendor/profile       - Get vendor profile
PUT    /api/v1/vendor/profile       - Update profile
GET    /api/v1/vendor/staff         - List staff
POST   /api/v1/vendor/staff         - Create staff member
POST   /api/v1/vendor/agreement/accept - Accept agreement
```

### Products (Vendor)
```
POST   /api/v1/vendor/products      - Create product (requires auth)
GET    /api/v1/vendor/products      - List own products
PUT    /api/v1/vendor/products/:id  - Update product
DELETE /api/v1/vendor/products/:id  - Archive product
POST   /api/v1/vendor/products/bulk-import - Bulk import via CSV
```

### Products (Public)
```
GET    /api/v1/products             - List published products
GET    /api/v1/products/:slug       - Get product details
GET    /api/v1/categories           - List categories
```

### Cart & Orders (Customer)
```
GET    /api/v1/customer/cart        - Get cart
POST   /api/v1/customer/cart/items  - Add to cart
PUT    /api/v1/customer/cart/items/:id - Update cart item
DELETE /api/v1/customer/cart/items/:id - Remove from cart
POST   /api/v1/customer/cart/apply-coupon - Apply discount code
POST   /api/v1/customer/orders/checkout - Place order
GET    /api/v1/customer/orders      - Order history
GET    /api/v1/customer/orders/:id  - Order details
POST   /api/v1/customer/returns     - Initiate return
```

### Orders (Vendor)
```
GET    /api/v1/vendor/orders        - View vendor orders
PUT    /api/v1/vendor/orders/:id/status - Update order status
```

### Admin
```
GET    /api/v1/admin/vendors        - List vendors
GET    /api/v1/admin/vendors/:id    - Vendor details
POST   /api/v1/admin/vendors/:id/approve - Approve vendor
POST   /api/v1/admin/vendors/:id/reject  - Reject vendor
POST   /api/v1/admin/vendors/:id/suspend - Suspend vendor
PUT    /api/v1/admin/vendors/:id/bank-details - Set bank details
PUT    /api/v1/admin/vendors/:id/commission  - Set commission

GET    /api/v1/admin/products/pending-approval - Pending products
POST   /api/v1/admin/products/:id/approve - Approve product
POST   /api/v1/admin/products/:id/reject  - Reject product

GET    /api/v1/admin/vendors/:id/wallet - Vendor wallet
POST   /api/v1/admin/vendors/:id/settle - Payout settlement
GET    /api/v1/admin/vendors/:id/settlements - Settlement history
```

---

## Database Models

### Users & Auth
- `users` — Platform users (admin, vendor, customer)
- `refresh_tokens` — Revokable refresh token hashes

### Vendors
- `vendors` — Vendor profiles with status (pending/approved/suspended)
- `vendor_bank_details` — Admin-controlled bank info
- `vendor_documents` — KYC document uploads
- `vendor_staff` — Vendor sub-accounts
- `vendor_wallet` — Balance tracking
- `vendor_agreements` — Versioned terms & conditions

### Catalog
- `categories` — Hierarchical product categories
- `brands` — Brand listings
- `products` — Product master with status workflow
- `product_variants` — SKU variants (size, color, etc.)
- `product_images` — Product images

### Orders
- `orders` — Master order (customer's unified cart)
- `sub_orders` — Vendor-specific orders (split from master)
- `order_items` — Individual items
- `order_status_logs` — Order status change history

### Returns & Refunds
- `return_requests` — Return/RMA requests
- `refunds` — Refund processing

### Financial
- `wallet_transactions` — Ledger entries
- `settlements` — Periodic payouts to vendors
- `promotions` — Discounts & coupons
- `audit_logs` — Admin action tracking

---

## Key Features Implemented

✅ **Authentication & Authorization**
- JWT access tokens (15 min expiry) + refresh tokens (7 days)
- Role-based access control (Super Admin, Admin, Vendor, Customer)
- Token revocation via hash storage

✅ **Vendor Onboarding**
- Structured registration workflow (pending → approved → active)
- KYC document upload to S3/MinIO
- Bank details (admin-only edit, vendor cannot change)
- Versioned vendor agreements with acceptance tracking

✅ **Product Catalog**
- Product status workflow (draft → pending_approval → published)
- Admin approval requirement before publication
- Variants, images, SEO metadata
- Category-level rules (min/max price, max discount %)

✅ **Order Processing**
- Single unified checkout across multiple vendors
- Automatic order splitting into sub-orders per vendor
- Stock validation before checkout
- Promotion priority: bank_ipg > product_discount > coupon

✅ **Commission & Settlement**
- Margin-based and markup-based commission models
- Commission priority: category > vendor > platform default
- Per-vendor wallet tracking
- Manual admin-initiated settlements

✅ **Database**
- PostgreSQL with GORM ORM
- Automatic migrations on startup
- Proper indexing for performance
- Soft deletes for audit trails

✅ **Caching**
- Redis for session cache, OTP storage, cart cache
- Generic cache helper with TTL support

✅ **File Storage**
- S3/MinIO integration for vendor documents & product images
- Upload with timestamp-based naming
- Delete support

---

## Next Steps: Full Implementation Checklist

### Phase 1: Core MVP (Current)
- [x] Project structure & scaffolding
- [x] Config & database connections
- [x] Models & migrations
- [x] JWT authentication
- [x] RBAC middleware
- [ ] **Vendor onboarding service** (register, upload docs, approve)
- [ ] **Product CRUD + approval** (create, list, approve, reject)
- [ ] **Cart & checkout** (add items, calculate total, split orders)
- [ ] **Order status management** (vendor updates, customer tracks)
- [ ] **Returns & refunds** (request, approve, calculate refund)
- [ ] **Settlement & payouts** (calculate vendor earnings, payout)

### Phase 2: Extended Features
- [ ] Promotions (CRUD, apply logic, usage tracking)
- [ ] Loyalty points (earn, redeem, expiry)
- [ ] WhatsApp integration (predefined templates, routing)
- [ ] Email notifications (order, return, settlement)
- [ ] Wishlist & reviews
- [ ] ERP integration (scheduled sync, error handling)

### Phase 3: Admin Dashboard
- [ ] Vendor management dashboard
- [ ] Product approval queue
- [ ] Order & return monitoring
- [ ] Commission & settlement reports
- [ ] Audit logs & compliance

---

## Testing

```bash
# Run unit tests (when implemented)
go test ./internal/services/... -v

# Run all tests with coverage
go test ./... -cover

# Integration tests with live DB (after docker-compose up)
go test -tags=integration ./internal/repositories/... -v
```

---

## Deployment

### Docker
```bash
# Build image
docker build -t coolmate-backend:latest .

# Run with docker-compose
docker-compose up
```

### Environment Variables
See `.env.example` for all configuration options.

---

## Development Notes

1. **Database Migrations**: Auto-run on server startup via `db.AutoMigrate()`
2. **Error Handling**: Services return `(result, error)` — handlers convert to HTTP responses
3. **Pagination**: Use `GetPaginationParams()` helper; default 20 items, max 100
4. **Logging**: Use `go.uber.org/zap` for structured logging (implement in Phase 2)
5. **Testing**: Use GORM transaction rollback for test isolation

---

## API Response Format

### Success Response
```json
{
  "success": true,
  "message": "Operation successful",
  "data": {
    "id": 1,
    "email": "user@example.com"
  }
}
```

### Paginated Response
```json
{
  "success": true,
  "message": "Data retrieved",
  "data": [...],
  "meta": {
    "total": 100,
    "page": 1,
    "pageSize": 20,
    "totalPages": 5
  }
}
```

### Error Response
```json
{
  "success": false,
  "message": "Bad Request",
  "errors": null
}
```

---

## Contact & Support

- **API Docs**: (Swagger endpoint TBD)
- **Issues**: (GitHub issues when repo is created)
- **Team**: Coolmate Team
