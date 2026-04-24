-- Initial database setup for Coolmate eCommerce
-- This script is run automatically by PostgreSQL initialization

-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Set up basic roles and permissions (optional, for security)
-- Note: Adjust user/password based on your .env configuration

-- Create initial schema
CREATE SCHEMA IF NOT EXISTS public;

-- The actual tables will be created by GORM AutoMigrate
-- This file serves as a placeholder for any manual SQL initialization
-- If you need custom SQL, add it below this line

-- Example: Create a custom index that GORM won't automatically create
-- CREATE INDEX idx_orders_vendor_id_status ON orders(vendor_id, status);
-- CREATE INDEX idx_products_vendor_id_status ON products(vendor_id, status);
