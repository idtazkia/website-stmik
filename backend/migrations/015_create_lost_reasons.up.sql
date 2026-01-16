-- Lost reasons for candidates who don't proceed
CREATE TABLE lost_reasons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    display_order INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index
CREATE INDEX idx_lost_reasons_is_active ON lost_reasons(is_active);
CREATE INDEX idx_lost_reasons_display_order ON lost_reasons(display_order);

-- Seed default lost reasons
INSERT INTO lost_reasons (name, description, display_order) VALUES
('Tidak ada respon', 'Kandidat tidak dapat dihubungi setelah beberapa kali percobaan', 1),
('Memilih kampus lain', 'Kandidat memutuskan mendaftar ke kampus kompetitor', 2),
('Kendala finansial', 'Kandidat tidak mampu membiayai kuliah', 3),
('Tidak memenuhi syarat', 'Kandidat tidak memenuhi persyaratan akademik atau administratif', 4),
('Waktu tidak tepat', 'Kandidat menunda kuliah ke tahun berikutnya', 5),
('Lokasi terlalu jauh', 'Kandidat keberatan dengan jarak tempuh ke kampus', 6),
('Orang tua tidak setuju', 'Keputusan orang tua untuk tidak melanjutkan', 7),
('Lainnya', 'Alasan lain yang tidak tercantum', 99);
