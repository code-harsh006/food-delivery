-- Initial database setup script
-- Run this script to set up the database schema

-- Create database (if not exists)
CREATE DATABASE IF NOT EXISTS servicehub;

-- Use the database
\c servicehub;

-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- The tables will be created automatically by GORM migrations
-- This script is for any additional setup needed

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);
CREATE INDEX IF NOT EXISTS idx_bookings_user_id ON bookings(user_id);
CREATE INDEX IF NOT EXISTS idx_bookings_service_id ON bookings(service_id);
CREATE INDEX IF NOT EXISTS idx_bookings_status ON bookings(status);
CREATE INDEX IF NOT EXISTS idx_bookings_scheduled_date ON bookings(scheduled_date);
CREATE INDEX IF NOT EXISTS idx_otps_user_id ON otps(user_id);
CREATE INDEX IF NOT EXISTS idx_otps_code ON otps(code);
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);

-- Insert initial data will be handled by the SeedServices function in Go
