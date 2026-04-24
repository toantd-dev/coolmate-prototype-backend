-- Setup script for Xiteb eCommerce PostgreSQL Database
-- Run this script on your PostgreSQL server before starting the application

-- Create database
CREATE DATABASE coolmate_ecommerce;

-- Connect to the newly created database and create extensions
\c coolmate_ecommerce;

-- Create extensions for UUID and crypto functions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- All tables will be created automatically by GORM AutoMigrate when the app starts
-- This script only needs to create the database and required extensions
