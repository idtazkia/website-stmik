-- Seed data for Phase 1

-- Test Users (for development/testing)
INSERT INTO users (email, name, role, is_active) VALUES
    ('admin@tazkia.ac.id', 'Admin User', 'admin', true),
    ('supervisor@tazkia.ac.id', 'Supervisor User', 'supervisor', true),
    ('consultant1@tazkia.ac.id', 'Consultant One', 'consultant', true),
    ('consultant2@tazkia.ac.id', 'Consultant Two', 'consultant', true);

-- Assign supervisor to consultants
UPDATE users SET id_supervisor = (SELECT id FROM users WHERE email = 'supervisor@tazkia.ac.id')
WHERE role = 'consultant';

-- Fee Types
INSERT INTO fee_types (name, code, is_recurring, installment_options) VALUES
    ('Biaya Pendaftaran', 'registration', false, '[1]'),
    ('SPP/Biaya Kuliah', 'tuition', true, '[1]'),
    ('Biaya Asrama', 'dormitory', true, '[1, 2, 10]');

-- Interaction Categories
INSERT INTO interaction_categories (name, sentiment, display_order) VALUES
    ('Tertarik', 'positive', 1),
    ('Mempertimbangkan', 'neutral', 2),
    ('Ragu-ragu', 'neutral', 3),
    ('Dingin', 'negative', 4),
    ('Tidak bisa dihubungi', 'negative', 5);

-- Obstacles
INSERT INTO obstacles (name, suggested_response, display_order) VALUES
    ('Biaya terlalu mahal', 'Jelaskan program beasiswa dan cicilan yang tersedia. Bandingkan dengan value yang didapat dari kuliah di STMIK Tazkia.', 1),
    ('Lokasi jauh', 'Jelaskan fasilitas asrama yang nyaman dan program pendampingan untuk mahasiswa dari luar kota.', 2),
    ('Orang tua belum setuju', 'Tawarkan untuk menghubungi orang tua langsung atau mengundang mereka untuk campus visit.', 3),
    ('Waktu belum tepat', 'Catat kapan waktu yang tepat dan jadwalkan follow-up. Informasikan tentang jadwal gelombang pendaftaran.', 4),
    ('Memilih kampus lain', 'Tanyakan kampus mana dan bandingkan keunggulan STMIK Tazkia. Catat sebagai feedback untuk improvement.', 5);

-- Default Prodis
INSERT INTO prodis (name, code, degree) VALUES
    ('Sistem Informasi', 'SI', 'S1'),
    ('Teknik Informatika', 'TI', 'S1');
