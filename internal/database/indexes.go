package database

import (
	"log"

	"github.com/coolmate/ecommerce-backend/internal/models"
	"gorm.io/gorm"
)

func CreateIndexes(db *gorm.DB) error {
	log.Println("Creating database indexes...")

	// Product indexes
	if !db.Migrator().HasIndex(&models.Product{}, "idx_vendor_id") {
		if err := db.Migrator().CreateIndex(&models.Product{}, "vendor_id"); err != nil {
			return err
		}
		log.Println("✓ Created index: products(vendor_id)")
	}

	if !db.Migrator().HasIndex(&models.Product{}, "idx_category_id") {
		if err := db.Migrator().CreateIndex(&models.Product{}, "category_id"); err != nil {
			return err
		}
		log.Println("✓ Created index: products(category_id)")
	}

	if !db.Migrator().HasIndex(&models.Product{}, "idx_status") {
		if err := db.Migrator().CreateIndex(&models.Product{}, "status"); err != nil {
			return err
		}
		log.Println("✓ Created index: products(status)")
	}

	if !db.Migrator().HasIndex(&models.Product{}, "idx_vendor_status") {
		if err := db.Migrator().CreateIndex(&models.Product{}, "idx_vendor_status"); err != nil {
			return err
		}
		log.Println("✓ Created index: products(vendor_id, status)")
	}

	// Order indexes
	if !db.Migrator().HasIndex(&models.Order{}, "idx_customer_id") {
		if err := db.Migrator().CreateIndex(&models.Order{}, "customer_id"); err != nil {
			return err
		}
		log.Println("✓ Created index: orders(customer_id)")
	}

	if !db.Migrator().HasIndex(&models.Order{}, "idx_order_status") {
		if err := db.Migrator().CreateIndex(&models.Order{}, "status"); err != nil {
			return err
		}
		log.Println("✓ Created index: orders(status)")
	}

	// SubOrder indexes
	if !db.Migrator().HasIndex(&models.SubOrder{}, "idx_suborder_vendor") {
		if err := db.Migrator().CreateIndex(&models.SubOrder{}, "vendor_id"); err != nil {
			return err
		}
		log.Println("✓ Created index: sub_orders(vendor_id)")
	}

	if !db.Migrator().HasIndex(&models.SubOrder{}, "idx_suborder_vendor_status") {
		if err := db.Migrator().CreateIndex(&models.SubOrder{}, "idx_suborder_vendor_status"); err != nil {
			return err
		}
		log.Println("✓ Created index: sub_orders(vendor_id, status)")
	}

	// Vendor indexes
	if !db.Migrator().HasIndex(&models.Vendor{}, "idx_vendor_status") {
		if err := db.Migrator().CreateIndex(&models.Vendor{}, "status"); err != nil {
			return err
		}
		log.Println("✓ Created index: vendors(status)")
	}

	if !db.Migrator().HasIndex(&models.Vendor{}, "idx_vendor_user") {
		if err := db.Migrator().CreateIndex(&models.Vendor{}, "user_id"); err != nil {
			return err
		}
		log.Println("✓ Created index: vendors(user_id)")
	}

	// User indexes (for faster lookups)
	if !db.Migrator().HasIndex(&models.User{}, "idx_user_email") {
		if err := db.Migrator().CreateIndex(&models.User{}, "email"); err != nil {
			return err
		}
		log.Println("✓ Created index: users(email)")
	}

	// VendorWallet indexes
	if !db.Migrator().HasIndex(&models.VendorWallet{}, "idx_wallet_vendor") {
		if err := db.Migrator().CreateIndex(&models.VendorWallet{}, "vendor_id"); err != nil {
			return err
		}
		log.Println("✓ Created index: vendor_wallets(vendor_id)")
	}

	// WalletTransaction indexes
	if !db.Migrator().HasIndex(&models.WalletTransaction{}, "idx_transaction_vendor") {
		if err := db.Migrator().CreateIndex(&models.WalletTransaction{}, "vendor_id"); err != nil {
			return err
		}
		log.Println("✓ Created index: wallet_transactions(vendor_id)")
	}

	// Settlement indexes
	if !db.Migrator().HasIndex(&models.Settlement{}, "idx_settlement_vendor") {
		if err := db.Migrator().CreateIndex(&models.Settlement{}, "vendor_id"); err != nil {
			return err
		}
		log.Println("✓ Created index: settlements(vendor_id)")
	}

	if !db.Migrator().HasIndex(&models.Settlement{}, "idx_settlement_status") {
		if err := db.Migrator().CreateIndex(&models.Settlement{}, "status"); err != nil {
			return err
		}
		log.Println("✓ Created index: settlements(status)")
	}

	// Category indexes
	if !db.Migrator().HasIndex(&models.Category{}, "idx_category_slug") {
		if err := db.Migrator().CreateIndex(&models.Category{}, "slug"); err != nil {
			return err
		}
		log.Println("✓ Created index: categories(slug)")
	}

	log.Println("✓ All database indexes created successfully")
	return nil
}
