USE iti_room_booking;

-- ==========================================================
-- 1. SEED USERS (5 Data)
-- Password Hash Valid untuk 'admin123'
-- Role disesuaikan dengan ENUM('admin', 'user')
-- ==========================================================
INSERT INTO users (name, role, email, password) VALUES 
('Super Admin', 'admin', 'admin@iti.ac.id', '$2a$10$vI8aWBnW3fID.ZQ4/zo1G.q1lRps.9cGLcZEiGDMVr5yUP1KUOYTa'),
('Dosen Pemdas', 'user', 'dosen@iti.ac.id', '$2a$10$vI8aWBnW3fID.ZQ4/zo1G.q1lRps.9cGLcZEiGDMVr5yUP1KUOYTa'),
('Kepala Lab', 'user', 'kalab@iti.ac.id', '$2a$10$vI8aWBnW3fID.ZQ4/zo1G.q1lRps.9cGLcZEiGDMVr5yUP1KUOYTa'),
('Admin Prodi', 'admin', 'prodi@iti.ac.id', '$2a$10$vI8aWBnW3fID.ZQ4/zo1G.q1lRps.9cGLcZEiGDMVr5yUP1KUOYTa'),
('Mahasiswa BEM', 'user', 'bem@iti.ac.id', '$2a$10$vI8aWBnW3fID.ZQ4/zo1G.q1lRps.9cGLcZEiGDMVr5yUP1KUOYTa');

-- ==========================================================
-- 2. SEED ROOMS (15 Data)
-- ==========================================================
INSERT INTO rooms (name, type, capacity, is_active) VALUES 
('Lab RPL', 'lab', 40, 1),
('Lab Jaringan', 'lab', 40, 1),
('Lab Multimedia', 'lab', 35, 1),
('Lab Hardware', 'lab', 30, 1),
('Lab AI & Data', 'lab', 40, 1),
('R. 201 (Teori)', 'kelas', 60, 1),
('R. 202 (Teori)', 'kelas', 60, 1),
('R. 203 (Teori)', 'kelas', 50, 1),
('R. 301 (Diskusi)', 'kelas', 30, 1),
('R. 302 (Diskusi)', 'kelas', 30, 1),
('Aula Serbaguna', 'aula', 200, 1),
('Auditorium Utama', 'aula', 500, 1),
('R. Sidang A', 'kelas', 20, 1),
('R. Sidang B', 'kelas', 20, 1),
('Gudang Arsip', 'gudang', 0, 0);

-- ==========================================================
-- 3. SEED BOOKINGS (Generator Otomatis 500 Data)
-- ==========================================================
DELIMITER $$

DROP PROCEDURE IF EXISTS GenerateDummyBookings$$

CREATE PROCEDURE GenerateDummyBookings()
BEGIN
    DECLARE i INT DEFAULT 1;
    DECLARE rand_room INT;
    DECLARE rand_user INT;
    DECLARE rand_month INT;
    DECLARE rand_day INT;
    DECLARE rand_hour INT;
    DECLARE rand_duration INT;
    DECLARE start_ts DATETIME;
    DECLARE end_ts DATETIME;
    DECLARE rand_status VARCHAR(20);
    DECLARE rand_purpose VARCHAR(100);
    
    WHILE i <= 500 DO
        -- Random IDs
        SET rand_room = FLOOR(1 + (RAND() * 15));   -- Room ID 1-15
        SET rand_user = FLOOR(1 + (RAND() * 5));    -- User ID 1-5 (ambil dari tabel users)
        
        -- Random Date & Time (Tahun 2026)
        SET rand_month = FLOOR(1 + (RAND() * 12));
        SET rand_day = FLOOR(1 + (RAND() * 28));    -- Aman untuk semua bulan
        SET rand_hour = FLOOR(8 + (RAND() * 9));    -- Jam 08:00 - 17:00
        SET rand_duration = FLOOR(1 + (RAND() * 3));-- Durasi 1-3 Jam
        
        -- Construct DATETIME
        SET start_ts = STR_TO_DATE(CONCAT('2026-', rand_month, '-', rand_day, ' ', rand_hour, ':00:00'), '%Y-%m-%d %H:%i:%s');
        SET end_ts = DATE_ADD(start_ts, INTERVAL rand_duration HOUR);
        
        -- Random Status (Sesuai ENUM)
        IF (RAND() < 0.7) THEN SET rand_status = 'approved';
        ELSEIF (RAND() < 0.9) THEN SET rand_status = 'pending';
        ELSE SET rand_status = 'rejected';
        END IF;

        -- Random Purpose
        IF (RAND() < 0.3) THEN SET rand_purpose = 'Kelas Tambahan';
        ELSEIF (RAND() < 0.6) THEN SET rand_purpose = 'Rapat Himpunan';
        ELSEIF (RAND() < 0.8) THEN SET rand_purpose = 'Sidang Skripsi';
        ELSE SET rand_purpose = 'Maintenance Rutin';
        END IF;

        -- Insert
        INSERT INTO bookings (room_id, user_id, start_time, end_time, purpose, status)
        VALUES (rand_room, rand_user, start_ts, end_ts, rand_purpose, rand_status);
        
        SET i = i + 1;
    END WHILE;
END$$

DELIMITER ;

-- Execute Generator
CALL GenerateDummyBookings();

-- Cleanup
DROP PROCEDURE IF EXISTS GenerateDummyBookings;

-- Cek Data
SELECT 'Setup & Seeding Complete' AS Status, 
        (SELECT COUNT(*) FROM users) AS Total_Users,
        (SELECT COUNT(*) FROM rooms) AS Total_Rooms,
        (SELECT COUNT(*) FROM bookings) AS Total_Bookings;