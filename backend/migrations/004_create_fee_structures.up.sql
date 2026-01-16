-- Fee structure per prodi and academic year
CREATE TABLE fee_structures (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_fee_type UUID NOT NULL REFERENCES fee_types(id) ON DELETE CASCADE,
    id_prodi UUID REFERENCES prodis(id) ON DELETE CASCADE,
    academic_year VARCHAR(9) NOT NULL,
    amount BIGINT NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(id_fee_type, id_prodi, academic_year)
);

-- Indexes
CREATE INDEX idx_fee_structures_id_fee_type ON fee_structures(id_fee_type);
CREATE INDEX idx_fee_structures_id_prodi ON fee_structures(id_prodi);
CREATE INDEX idx_fee_structures_academic_year ON fee_structures(academic_year);
CREATE INDEX idx_fee_structures_is_active ON fee_structures(is_active);

-- Trigger for updated_at
CREATE TRIGGER trigger_fee_structures_updated_at
    BEFORE UPDATE ON fee_structures
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
