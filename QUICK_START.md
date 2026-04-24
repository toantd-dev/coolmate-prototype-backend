# Quick Start Guide

## 1. Start Local Services

```bash
# Start PostgreSQL, Redis, MinIO
docker-compose up -d

# Verify containers are running
docker-compose ps
```

Expected output:
```
NAME           STATUS
coolmate_postgres Running
coolmate_redis    Running
coolmate_minio    Running
```

## 2. Install Dependencies & Run Server

```bash
# Download Go dependencies
go mod tidy

# Start the API server
go run cmd/api/main.go
```

Expected output:
```
✓ Connected to PostgreSQL successfully
✓ Connected to Redis successfully
[GIN-debug] Loaded HTML Templates
[GIN-debug] listening on .:8080
```

## 3. Test Authentication Endpoints

### Register a Customer

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "SecurePassword123",
    "firstName": "John",
    "lastName": "Doe",
    "phone": "+1234567890"
  }'
```

Response:
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "email": "john@example.com",
      "firstName": "John",
      "lastName": "Doe",
      "role": "customer"
    }
  }
}
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "SecurePassword123"
  }'
```

### Refresh Token

```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }'
```

### Logout

```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Content-Type: application/json" \
  -d '{
    "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }'
```

## 4. Test Protected Endpoints

Use the `accessToken` from registration/login for protected routes:

```bash
# Example: Get vendor profile (requires auth)
curl -X GET http://localhost:8080/api/v1/vendor/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

## 5. Using Postman/Insomnia

**Create a collection** with these requests:

1. **Register Customer** - POST `/api/v1/auth/register`
   - Body: JSON with email, password, firstName, lastName, phone
   - Save response token in Postman environment: `{{accessToken}}`

2. **Login** - POST `/api/v1/auth/login`
   - Body: JSON with email, password

3. **Refresh Token** - POST `/api/v1/auth/refresh`
   - Body: JSON with refreshToken

4. **Get Vendor Profile** - GET `/api/v1/vendor/profile`
   - Header: `Authorization: Bearer {{accessToken}}`
   - (Currently returns stub response)

## 6. Access MinIO Console

Open browser: `http://localhost:9001`
- Username: `minioadmin`
- Password: `minioadmin`

Create a bucket named `coolmate` for file uploads.

## 7. PostgreSQL Access

```bash
# Connect to database
psql postgresql://coolmate_user:coolmate_password@localhost:5432/coolmate_ecommerce

# List tables
\dt

# View users
SELECT id, email, role, status FROM users;
```

## 8. Redis Access

```bash
# Connect to Redis
redis-cli -h localhost -p 6379

# Check connected
PING
# Should return: PONG

# View keys (during testing)
KEYS *
```

## Common Issues

### Port Already in Use

```bash
# Kill process on port 8080
lsof -ti:8080 | xargs kill -9

# Or use different port in .env
SERVER_PORT=8081
```

### PostgreSQL Connection Failed

```bash
# Check if container is running
docker-compose logs postgres

# Restart containers
docker-compose restart postgres
```

### "Invalid JWT Secret" Error

Make sure `.env` has a valid `JWT_SECRET`:
```
JWT_SECRET=your_super_secret_key_at_least_32_chars_long
```

### S3 Upload Fails

1. Create bucket in MinIO console
2. Ensure `S3_BUCKET=coolmate` in `.env`
3. Check MinIO logs: `docker-compose logs minio`

## Next: Test Full Workflows

Once auth is working, test these workflows as they're implemented:

1. **Vendor Onboarding** → Register vendor → Upload docs → Admin approve
2. **Product Creation** → Vendor creates product → Admin approves → Visible in catalog
3. **Shopping Cart** → Add items → Apply coupon → Checkout
4. **Order Management** → Vendor updates status → Customer receives → Return flow

---

## Development Tips

### Hot Reload
```bash
# Install air for hot reload
go install github.com/cosmtrek/air@latest

# Run with air
air
```

### View Logs
```bash
# PostgreSQL logs
docker-compose logs -f postgres

# Redis logs
docker-compose logs -f redis

# MinIO logs
docker-compose logs -f minio
```

### Reset Database
```bash
# Stop containers and remove volumes
docker-compose down -v

# Restart with fresh DB
docker-compose up -d
```

### Debug Mode
Set in `.env`:
```
GIN_MODE=debug
LOG_LEVEL=debug
```

---

## Environment Variables

Copy from `.env.example` and customize:

```bash
cp .env.example .env
# Edit .env with your settings
```

Key variables:
- `SERVER_PORT` — API port (default 8080)
- `DB_*` — PostgreSQL connection
- `REDIS_*` — Redis connection
- `JWT_SECRET` — Must be 32+ chars
- `S3_*` — MinIO settings
- `CORS_ORIGINS` — Allowed origins

---

## What's Ready to Use

✅ **Authentication** — Register, login, refresh, logout working
✅ **Authorization** — Role-based access control middleware ready
✅ **Database** — All tables auto-created on startup
✅ **File Upload** — S3/MinIO integration ready
✅ **Caching** — Redis helpers for session/cache
✅ **API Structure** — Routes organized, response formatting consistent

## What Needs Implementation

⏳ **Vendor Services** — Register vendor, upload docs, list, manage
⏳ **Product Services** — CRUD, approval workflow, inventory
⏳ **Order Services** — Cart, checkout, order splitting, commission calc
⏳ **Return/Refund** — Processing, wallet updates
⏳ **Settlement** — Vendor payouts

See [IMPLEMENTATION_STATUS.md](IMPLEMENTATION_STATUS.md) for detailed checklist.

---

## Get Help

1. Check the code comments in `internal/models/` for schema details
2. Read `README.md` for API documentation
3. Review `IMPLEMENTATION_STATUS.md` for what's next
4. Check CloudFormation logs: `docker-compose logs <service>`
