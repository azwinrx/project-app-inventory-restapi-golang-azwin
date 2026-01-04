-- =============================================
-- INVENTORY MANAGEMENT SYSTEM DATABASE
-- =============================================

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- =============================================
-- DROP TABLES (untuk reset jika sudah ada)
-- =============================================
DROP TABLE IF EXISTS sale_items CASCADE;
DROP TABLE IF EXISTS sales CASCADE;
DROP TABLE IF EXISTS items CASCADE;
DROP TABLE IF EXISTS racks CASCADE;
DROP TABLE IF EXISTS warehouses CASCADE;
DROP TABLE IF EXISTS categories CASCADE;
DROP TABLE IF EXISTS sessions CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- =============================================
-- CREATE TABLES
-- =============================================

-- Table: users
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('super_admin', 'admin', 'staff')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table: sessions (untuk authentication token)
CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token UUID NOT NULL DEFAULT uuid_generate_v4(),
    expired_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table: categories
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table: warehouses
CREATE TABLE warehouses (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    location VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table: racks
CREATE TABLE racks (
    id SERIAL PRIMARY KEY,
    warehouse_id INTEGER NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table: items (barang)
CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    rack_id INTEGER NOT NULL REFERENCES racks(id) ON DELETE RESTRICT,
    name VARCHAR(150) NOT NULL,
    sku VARCHAR(50) UNIQUE NOT NULL,
    stock INTEGER NOT NULL DEFAULT 0,
    min_stock INTEGER NOT NULL DEFAULT 5,
    price DECIMAL(15, 2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table: sales (transaksi penjualan)
CREATE TABLE sales (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    total_amount DECIMAL(15, 2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table: sale_items (detail item penjualan)
CREATE TABLE sale_items (
    id SERIAL PRIMARY KEY,
    sale_id INTEGER NOT NULL REFERENCES sales(id) ON DELETE CASCADE,
    item_id INTEGER NOT NULL REFERENCES items(id) ON DELETE RESTRICT,
    quantity INTEGER NOT NULL,
    price DECIMAL(15, 2) NOT NULL,
    subtotal DECIMAL(15, 2) NOT NULL
);

-- =============================================
-- CREATE INDEXES
-- =============================================
CREATE INDEX idx_sessions_token ON sessions(token);
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_items_category_id ON items(category_id);
CREATE INDEX idx_items_rack_id ON items(rack_id);
CREATE INDEX idx_items_stock ON items(stock);
CREATE INDEX idx_racks_warehouse_id ON racks(warehouse_id);
CREATE INDEX idx_sales_user_id ON sales(user_id);
CREATE INDEX idx_sale_items_sale_id ON sale_items(sale_id);
CREATE INDEX idx_sale_items_item_id ON sale_items(item_id);

-- =============================================
-- INSERT SAMPLE DATA
-- =============================================

-- Users (password: password123 - hashed dengan bcrypt)
-- Hash: $2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqOqWL3QRPaLIvMvmCdq8DhKxJEau
INSERT INTO users (username, email, password, role) VALUES
('superadmin', 'superadmin@inventory.com', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqOqWL3QRPaLIvMvmCdq8DhKxJEau', 'super_admin'),
('admin1', 'admin1@inventory.com', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqOqWL3QRPaLIvMvmCdq8DhKxJEau', 'admin'),
('staff1', 'staff1@inventory.com', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqOqWL3QRPaLIvMvmCdq8DhKxJEau', 'staff'),
('staff2', 'staff2@inventory.com', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqOqWL3QRPaLIvMvmCdq8DhKxJEau', 'staff');

-- Categories
INSERT INTO categories (name) VALUES
('Elektronik'),
('Pakaian'),
('Makanan & Minuman'),
('Alat Tulis'),
('Perabotan');

-- Warehouses
INSERT INTO warehouses (name, location) VALUES
('Gudang Utama', 'Jl. Industri No. 1, Jakarta'),
('Gudang Cabang', 'Jl. Raya No. 25, Bandung'),
('Gudang Timur', 'Jl. Pelabuhan No. 10, Surabaya');

-- Racks
INSERT INTO racks (warehouse_id, name) VALUES
(1, 'Rak A1'),
(1, 'Rak A2'),
(1, 'Rak B1'),
(2, 'Rak C1'),
(2, 'Rak C2'),
(3, 'Rak D1');

-- Items
INSERT INTO items (category_id, rack_id, name, sku, stock, min_stock, price) VALUES
(1, 1, 'Laptop Asus ROG', 'ELK-001', 15, 5, 15000000.00),
(1, 1, 'Mouse Logitech', 'ELK-002', 50, 10, 350000.00),
(1, 2, 'Keyboard Mechanical', 'ELK-003', 3, 5, 750000.00),
(2, 3, 'Kaos Polos Hitam', 'PKN-001', 100, 20, 75000.00),
(2, 3, 'Celana Jeans', 'PKN-002', 4, 5, 250000.00),
(3, 4, 'Kopi Arabica 250gr', 'MKN-001', 200, 30, 85000.00),
(3, 4, 'Teh Hijau 100gr', 'MKN-002', 2, 5, 45000.00),
(4, 5, 'Pulpen Pilot', 'ATK-001', 500, 50, 5000.00),
(4, 5, 'Buku Tulis A5', 'ATK-002', 300, 50, 8000.00),
(5, 6, 'Meja Kantor', 'PRB-001', 10, 5, 1500000.00),
(5, 6, 'Kursi Ergonomis', 'PRB-002', 1, 5, 2500000.00);

-- Sales
INSERT INTO sales (user_id, total_amount, created_at) VALUES
(3, 15350000.00, '2025-12-15 10:30:00'),
(3, 500000.00, '2025-12-18 14:20:00'),
(4, 170000.00, '2025-12-20 09:15:00'),
(3, 2585000.00, '2025-12-25 11:45:00'),
(4, 40000.00, '2025-12-28 16:00:00');

-- Sale Items
INSERT INTO sale_items (sale_id, item_id, quantity, price, subtotal) VALUES
-- Sale 1: Laptop + Mouse
(1, 1, 1, 15000000.00, 15000000.00),
(1, 2, 1, 350000.00, 350000.00),
-- Sale 2: Celana Jeans x2
(2, 5, 2, 250000.00, 500000.00),
-- Sale 3: Kopi Arabica x2
(3, 6, 2, 85000.00, 170000.00),
-- Sale 4: Keyboard + Kursi
(4, 3, 1, 750000.00, 750000.00),
(4, 11, 1, 2500000.00, 2500000.00),
-- Sale 5: Pulpen x5 + Buku x3 (untuk update stock jadi dibawah minimum untuk testing)
(5, 8, 5, 5000.00, 25000.00);

-- Update stock setelah penjualan (simulasi)
UPDATE items SET stock = stock - 1 WHERE id = 1; -- Laptop: 15 -> 14
UPDATE items SET stock = stock - 1 WHERE id = 2; -- Mouse: 50 -> 49
UPDATE items SET stock = stock - 2 WHERE id = 5; -- Celana: 4 -> 2 (dibawah min)
UPDATE items SET stock = stock - 2 WHERE id = 6; -- Kopi: 200 -> 198
UPDATE items SET stock = stock - 1 WHERE id = 3; -- Keyboard: 3 -> 2 (dibawah min)
UPDATE items SET stock = stock - 1 WHERE id = 11; -- Kursi: 1 -> 0 (dibawah min)
UPDATE items SET stock = stock - 5 WHERE id = 8; -- Pulpen: 500 -> 495

-- =============================================
-- VERIFICATION QUERIES
-- =============================================

-- Lihat semua users
-- SELECT * FROM users;

-- Lihat items dengan stock dibawah minimum (untuk fitur cek stok minimum)
-- SELECT i.*, c.name as category_name, r.name as rack_name, w.name as warehouse_name
-- FROM items i
-- JOIN categories c ON i.category_id = c.id
-- JOIN racks r ON i.rack_id = r.id
-- JOIN warehouses w ON r.warehouse_id = w.id
-- WHERE i.stock < i.min_stock;

-- Lihat report penjualan
-- SELECT 
--     COUNT(DISTINCT s.id) as total_sales,
--     SUM(s.total_amount) as total_revenue,
--     COUNT(DISTINCT si.item_id) as total_items_sold,
--     SUM(si.quantity) as total_quantity_sold
-- FROM sales s
-- JOIN sale_items si ON s.id = si.sale_id;
