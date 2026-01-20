CREATE DATABASE IF NOT EXISTS iti_room_booking;
USE iti_room_booking;

-- Matikan Foreign Key Check agar bisa drop table dengan aman
SET FOREIGN_KEY_CHECKS = 0;

-- 1. Reset Tables
DROP TABLE IF EXISTS bookings;
DROP TABLE IF EXISTS rooms;
DROP TABLE IF EXISTS users;

-- 2. Table Users
-- Sesuai gambar: role adalah ENUM('admin', 'user')
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    role ENUM('admin', 'user') NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL
    -- created_at opsional, bisa ditambahkan jika perlu history akun
);

-- 3. Table Rooms
-- Sesuai gambar: type adalah VARCHAR(50), is_active tinyint(1)
CREATE TABLE rooms (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL, -- Flexible ('lab', 'kelas', 'aula')
    capacity INT NOT NULL,
    is_active TINYINT(1) DEFAULT 1
);

-- 4. Table Bookings
-- Sesuai gambar: start_time & end_time DATETIME, status ENUM
CREATE TABLE bookings (
    id INT AUTO_INCREMENT PRIMARY KEY,
    room_id INT NOT NULL,
    user_id INT NOT NULL,
    start_time DATETIME NOT NULL,
    end_time DATETIME NOT NULL,
    purpose TEXT,
    status ENUM('pending', 'approved', 'rejected') DEFAULT 'pending',
    
    -- Foreign Keys (Menjaga integritas data)
    CONSTRAINT fk_booking_room FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,
    CONSTRAINT fk_booking_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Nyalakan kembali Foreign Key Check
SET FOREIGN_KEY_CHECKS = 1;