# Setup Local Development Environment

## Prerequisites
- Go 1.22+
- Docker & Docker Compose
- PostgreSQL client (optional, for manual setup)
- PowerShell 7+ or CMD

## Quick Start (Recommended)

### 1. Start Docker Services
```powershell
# From project root directory
docker-compose up -d

# Verify services are running
docker-compose ps
```

This starts:
- PostgreSQL on localhost:5432
- Redis on localhost:6379
- MinIO on localhost:9000

### 2. Create .env File
```powershell
# Copy from template (if exists) or create manually
$env_content = @"
# Server
SERVER_PORT=8080
SERVER_MODE=debug

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=coolmate_ecommerce

# Redis
REDIS_URL=redis://localhost:6379/0

# S3 / MinIO
S3_ENDPOINT=http://localhost:9000
S3_REGION=us-east-1
S3_ACCESS_KEY=minioadmin
S3_SECRET_KEY=minioadmin
S3_BUCKET=ecommerce

# JWT
JWT_SECRET=your-super-secret-key-change-this-in-production
JWT_ACCESS_TOKEN_EXPIRE_MINUTES=60
JWT_REFRESH_TOKEN_EXPIRE_DAYS=7

# CORS
CORS_ORIGINS=http://localhost:3000,http://localhost:8080

# Email (Optional - for production)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
"@

$env_content | Out-File -FilePath .env -Encoding UTF8
```

### 3. Install Go Dependencies
```powershell
go mod download
go mod tidy
```

### 4. Run the Application
```powershell
# Database tables will be created automatically via AutoMigrate
go run cmd/api/main.go
```

**Expected output:**
```
Starting API server on :8080
Metrics server running on :9090
```

### 5. Test the API
```powershell
# Health check endpoint
curl http://localhost:8080/health

# Should return:
# {"status":"healthy"}
```

---

## Manual Database Setup (If Not Using Docker)

### Option A: PostgreSQL Local Installation
```powershell
# Run the setup script
psql -U postgres -h localhost -f setup-db.sql
```

### Option B: Manual Commands
```powershell
# Connect to PostgreSQL
psql -U postgres -h localhost

# Then run these commands:
# CREATE DATABASE coolmate_ecommerce;
# \c coolmate_ecommerce;
# CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
# CREATE EXTENSION IF NOT EXISTS "pgcrypto";
# \q
```

---

## Project Structure
```
coolmate-prototype backend/
├── cmd/api/main.go              # Entry point - handles AutoMigrate
├── internal/
│   ├── database/
│   │   ├── postgres.go          # PostgreSQL connection
│   │   ├── indexes.go           # Custom database indexes
│   │   └── redis.go             # Redis connection
│   ├── models/                  # GORM models (auto-migrated)
│   ├── handlers/                # HTTP handlers
│   ├── services/                # Business logic
│   ├── repositories/            # Data access layer
│   ├── middleware/              # Auth & CORS middleware
│   └── config/                  # Configuration loading
├── pkg/
│   ├── auth/                    # JWT utilities
│   ├── cache/                   # Redis cache wrapper
│   └── storage/                 # S3 upload utilities
├── migrations/
│   └── init-db.sql             # Database initialization
├── docker-compose.yml           # Local development services
├── .env                        # Environment variables (create this)
└── go.mod                      # Go module definition
```

---

## Common Issues & Solutions

### Issue: "Failed to connect to database"
**Solution:**
- Check PostgreSQL is running: `docker-compose ps`
- Verify `.env` file has correct DB credentials
- Ensure database `coolmate_ecommerce` exists

### Issue: "Failed to connect to Redis"
**Solution:**
- Check Redis is running: `docker-compose ps`
- Verify Redis URL in `.env`: `redis://localhost:6379/0`

### Issue: "Failed to initialize S3"
**Solution:**
- Check MinIO is running: `docker-compose ps`
- MinIO credentials in `.env` should be:
  - `S3_ACCESS_KEY=minioadmin`
  - `S3_SECRET_KEY=minioadmin`

### Issue: AutoMigrate Fails
**Solution:**
- Database must exist first (created by setup-db.sql)
- Check CreateIndexes function in `internal/database/indexes.go`
- Review GORM model definitions in `internal/models/`

---

## Useful Commands

```powershell
# Run tests with coverage
go test ./... -cover

# Format code
go fmt ./...

# Lint code
go vet ./...

# View metrics
curl http://localhost:9090/metrics

# Stop all Docker services
docker-compose down

# View logs
docker-compose logs -f postgres
docker-compose logs -f redis
docker-compose logs -f minio
```

---

## API Endpoints

After starting the server, key endpoints are:

**Health Check:**
```
GET http://localhost:8080/health
```

**Authentication:**
```
POST http://localhost:8080/api/v1/auth/register
POST http://localhost:8080/api/v1/auth/login
```

**For more endpoints, see setupRoutes() in cmd/api/main.go**

---

## Database Auto-Migration Details

The application automatically creates all tables on startup via GORM's `AutoMigrate()`:

✅ Users, Vendors, Products, Orders, etc.  
✅ Foreign key relationships  
✅ Indexes (via `CreateIndexes()` function)  
✅ No manual SQL needed!

---

## Next Steps

1. ✅ Start Docker services: `docker-compose up -d`
2. ✅ Create `.env` file with your config
3. ✅ Run: `go run cmd/api/main.go`
4. ✅ Test: `curl http://localhost:8080/health`
5. 📚 Check API documentation in `main.go` for endpoint details

**Happy coding! 🚀**
