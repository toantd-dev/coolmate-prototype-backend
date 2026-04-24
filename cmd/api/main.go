package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/coolmate/ecommerce-backend/internal/config"
	"github.com/coolmate/ecommerce-backend/internal/database"
	_ "github.com/coolmate/ecommerce-backend/docs"
	"github.com/coolmate/ecommerce-backend/internal/middleware"
	"github.com/coolmate/ecommerce-backend/internal/models"
	"github.com/coolmate/ecommerce-backend/internal/handlers"
	"github.com/coolmate/ecommerce-backend/internal/repositories"
	"github.com/coolmate/ecommerce-backend/internal/services"
	"github.com/coolmate/ecommerce-backend/pkg/auth"
	"github.com/coolmate/ecommerce-backend/pkg/cache"
	"github.com/coolmate/ecommerce-backend/pkg/jsondb"
	"github.com/coolmate/ecommerce-backend/pkg/storage"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Coolmate Multivendor eCommerce API
// @version         0.1
// @description     Go + Gin backend for a multivendor eCommerce platform. Auth, vendor onboarding (KYC), product catalog, order splitting, commission, settlement.
// @termsOfService  https://example.com/terms
// @contact.name    Coolmate Team
// @contact.email   dev@coolmate.local
// @license.name    MIT
// @host            localhost:8080
// @BasePath        /api/v1
// @schemes         http https
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
// @description     Paste: `Bearer <accessToken>`

func main() {
	// Load config from .env
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set gin mode
	gin.SetMode(cfg.Server.Mode)

	// Connect to database
	db, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate models
	if err := db.AutoMigrate(
		&models.User{},
		&models.RefreshToken{},
		&models.Vendor{},
		&models.VendorBankDetails{},
		&models.VendorDocument{},
		&models.VendorStaff{},
		&models.VendorWallet{},
		&models.VendorAgreement{},
		&models.VendorAgreementAcceptance{},
		&models.Category{},
		&models.Brand{},
		&models.Product{},
		&models.ProductVariant{},
		&models.ProductImage{},
		&models.ProductApprovalLog{},
		&models.Cart{},
		&models.CartItem{},
		&models.Order{},
		&models.SubOrder{},
		&models.OrderItem{},
		&models.OrderStatusLog{},
		&models.ReturnRequest{},
		&models.Refund{},
		&models.Promotion{},
		&models.ProductPromotion{},
		&models.WalletTransaction{},
		&models.Settlement{},
		&models.AuditLog{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Create database indexes for optimal query performance.
	// Auto-migrate already creates struct-tag indexes; non-fatal if extras fail.
	if err := database.CreateIndexes(db); err != nil {
		log.Printf("Warning: optional custom indexes skipped: %v", err)
	}

	// Connect to Redis
	redisClient, err := database.NewRedisClient(&cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Initialize S3 manager
	s3Manager, err := storage.NewS3Manager(&cfg.S3)
	if err != nil {
		log.Fatalf("Failed to initialize S3: %v", err)
	}

	// Initialize cache manager
	cacheManager := cache.NewCacheManager(redisClient)

	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenExpireMinutes,
		cfg.JWT.RefreshTokenExpireDays,
	)

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	vendorRepo := repositories.NewVendorRepository(db)

	// Use JSON repository for products if USE_JSON_DATA env var is set
	var productRepo repositories.IProductRepository
	if cfg.Server.Mode == "debug" {
		// Try to load from JSON in debug mode
		jsonLoader, err := initJSONLoader()
		if err == nil {
			log.Println("Using JSON data for products")
			productRepo = jsondb.NewJSONProductRepository(jsonLoader)
		} else {
			log.Println("JSON data not found, using database for products")
			productRepo = repositories.NewProductRepository(db)
		}
	} else {
		productRepo = repositories.NewProductRepository(db)
	}

	orderRepo := repositories.NewOrderRepository(db)
	promotionRepo := repositories.NewPromotionRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, jwtManager, cacheManager)
	vendorService := services.NewVendorService(vendorRepo, userRepo, s3Manager, cacheManager)
	productService := services.NewProductService(productRepo, vendorRepo, cacheManager)
	orderService := services.NewOrderService(orderRepo, productRepo, promotionRepo, vendorRepo, cacheManager)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	vendorHandler := handlers.NewVendorHandler(vendorService)
	productHandler := handlers.NewProductHandler(productService)
	orderHandler := handlers.NewOrderHandler(orderService)
	healthHandler := handlers.NewHealthHandler(db)

	// Create router
	router := gin.Default()

	// CORS
	router.Use(middleware.CORS(cfg.CORS.Origins))

	// Setup routes
	setupRoutes(router, authHandler, vendorHandler, productHandler, orderHandler, healthHandler, jwtManager)

	// Start metrics server on port 9090
	metricsHandler := handlers.NewMetricsHandler()
	go func() {
		metricsServer := &http.Server{
			Addr: ":9090",
		}
		http.HandleFunc("/metrics", metricsHandler.ServeMetrics)
		log.Println("Starting metrics server on :9090")
		if err := metricsServer.ListenAndServe(); err != nil {
			log.Printf("Metrics server error: %v", err)
		}
	}()

	// Start API server on configured port
	log.Printf("Starting API server on :%s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupRoutes(r *gin.Engine, authHandler *handlers.AuthHandler, vendorHandler *handlers.VendorHandler, productHandler *handlers.ProductHandler, orderHandler *handlers.OrderHandler, healthHandler *handlers.HealthHandler, jwtManager *auth.JWTManager) {
	// Health check endpoint (not under /api/v1)
	r.GET("/health", healthHandler.HealthCheck)

	// Swagger UI — browse at /swagger/index.html
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")

	// Auth routes (public)
	authRoutes := v1.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/refresh", authHandler.Refresh)
		authRoutes.POST("/logout", authHandler.Logout)
	}

	// Protected routes
	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware(jwtManager))

	// Vendor routes
	vendorRoutes := protected.Group("/vendor")
	{
		vendorRoutes.POST("/register", authHandler.RegisterVendor)
		vendorRoutes.POST("/documents/upload", vendorHandler.UploadDocument)
		vendorRoutes.GET("/profile", vendorHandler.GetProfile)
		vendorRoutes.PUT("/profile", vendorHandler.UpdateProfile)
		vendorRoutes.GET("/staff", vendorHandler.GetStaff)
		vendorRoutes.POST("/staff", vendorHandler.CreateStaff)
		vendorRoutes.PUT("/staff/:id", vendorHandler.UpdateStaff)
		vendorRoutes.DELETE("/staff/:id", vendorHandler.DeleteStaff)
		vendorRoutes.POST("/agreement/accept", vendorHandler.AcceptAgreement)

		// Product routes for vendors
		vendorRoutes.POST("/products", productHandler.CreateProduct)
		vendorRoutes.GET("/products", productHandler.ListVendorProducts)
		vendorRoutes.PUT("/products/:id", productHandler.UpdateProduct)
		vendorRoutes.DELETE("/products/:id", productHandler.ArchiveProduct)
		vendorRoutes.POST("/products/bulk-import", productHandler.BulkImportProducts)

		// Order routes for vendors
		vendorRoutes.GET("/orders", orderHandler.ListVendorOrders)
		vendorRoutes.PUT("/orders/:id/status", orderHandler.UpdateOrderStatus)
	}

	// Customer routes
	customerRoutes := protected.Group("/customer")
	{
		// Product browse
		customerRoutes.GET("/products", productHandler.ListProducts)
		customerRoutes.GET("/products/:slug", productHandler.GetProductBySlug)

		// Cart
		customerRoutes.GET("/cart", orderHandler.GetCart)
		customerRoutes.POST("/cart/items", orderHandler.AddToCart)
		customerRoutes.PUT("/cart/items/:id", orderHandler.UpdateCartItem)
		customerRoutes.DELETE("/cart/items/:id", orderHandler.RemoveFromCart)
		customerRoutes.POST("/cart/apply-coupon", orderHandler.ApplyCoupon)

		// Orders
		customerRoutes.POST("/orders/checkout", orderHandler.Checkout)
		customerRoutes.GET("/orders", orderHandler.ListCustomerOrders)
		customerRoutes.GET("/orders/:id", orderHandler.GetOrderDetail)

		// Returns
		customerRoutes.POST("/returns", orderHandler.InitiateReturn)
		customerRoutes.GET("/returns", orderHandler.ListReturns)
	}

	// Admin routes
	adminRoutes := protected.Group("/admin")
	adminRoutes.Use(middleware.RequireAdmin())
	{
		// Vendor management
		adminRoutes.GET("/vendors", vendorHandler.ListVendors)
		adminRoutes.GET("/vendors/:id", vendorHandler.GetVendor)
		adminRoutes.POST("/vendors/:id/approve", vendorHandler.ApproveVendor)
		adminRoutes.POST("/vendors/:id/reject", vendorHandler.RejectVendor)
		adminRoutes.POST("/vendors/:id/suspend", vendorHandler.SuspendVendor)
		adminRoutes.PUT("/vendors/:id/bank-details", vendorHandler.UpdateBankDetails)
		adminRoutes.PUT("/vendors/:id/commission", vendorHandler.UpdateCommission)

		// Product approval
		adminRoutes.GET("/products/pending-approval", productHandler.ListPendingProducts)
		adminRoutes.POST("/products/:id/approve", productHandler.ApproveProduct)
		adminRoutes.POST("/products/:id/reject", productHandler.RejectProduct)

		// Orders
		adminRoutes.GET("/orders", orderHandler.ListAllOrders)

		// Settlements
		adminRoutes.GET("/vendors/:id/wallet", vendorHandler.GetVendorWallet)
		adminRoutes.POST("/vendors/:id/settle", vendorHandler.SettleVendor)
		adminRoutes.GET("/vendors/:id/settlements", vendorHandler.GetSettlements)
	}

	// Public routes
	publicRoutes := v1.Group("")
	{
		publicRoutes.GET("/products", productHandler.ListProducts)
		publicRoutes.GET("/products/:slug", productHandler.GetProductBySlug)
		publicRoutes.GET("/categories", productHandler.GetCategories)
	}
}

func initJSONLoader() (*jsondb.JSONLoader, error) {
	return jsondb.NewJSONLoader("pkg/jsondb/products.json")
}
