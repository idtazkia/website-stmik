-- Document Types: Configure required documents for candidates
CREATE TABLE document_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    code VARCHAR(30) NOT NULL UNIQUE,
    description TEXT,
    is_required BOOLEAN DEFAULT true,
    can_defer BOOLEAN DEFAULT false,
    max_file_size_mb INT DEFAULT 5,
    allowed_extensions TEXT[] DEFAULT ARRAY['jpg', 'jpeg', 'png', 'pdf'],
    display_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_document_types_is_active ON document_types(is_active);
CREATE INDEX idx_document_types_display_order ON document_types(display_order);

-- Seed default document types
INSERT INTO document_types (name, code, description, is_required, can_defer, max_file_size_mb, display_order) VALUES
    ('KTP/Kartu Identitas', 'ktp', 'Kartu Tanda Penduduk atau kartu identitas lainnya', true, false, 2, 1),
    ('Pas Foto', 'photo', 'Pas foto 3x4 latar belakang merah', true, false, 2, 2),
    ('Ijazah SMA/SMK', 'ijazah', 'Ijazah atau Surat Keterangan Lulus', true, true, 5, 3),
    ('Transkrip Nilai', 'transcript', 'Transkrip nilai atau rapor semester terakhir', true, true, 5, 4);
