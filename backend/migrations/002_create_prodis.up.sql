-- Program Studi (Prodi) table
CREATE TABLE prodis (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    code VARCHAR(10) NOT NULL UNIQUE,
    degree VARCHAR(10) NOT NULL CHECK (degree IN ('S1', 'D3')),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_prodis_code ON prodis(code);
CREATE INDEX idx_prodis_is_active ON prodis(is_active);

-- Trigger for updated_at
CREATE TRIGGER trigger_prodis_updated_at
    BEFORE UPDATE ON prodis
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
