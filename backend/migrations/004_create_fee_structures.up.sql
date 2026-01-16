-- Fee structure per prodi and academic year
CREATE TABLE fee_structures (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fee_type_id UUID NOT NULL REFERENCES fee_types(id) ON DELETE CASCADE,
    prodi_id UUID REFERENCES prodis(id) ON DELETE CASCADE,
    academic_year VARCHAR(9) NOT NULL,
    amount BIGINT NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(fee_type_id, prodi_id, academic_year)
);

-- Indexes
CREATE INDEX idx_fee_structures_fee_type_id ON fee_structures(fee_type_id);
CREATE INDEX idx_fee_structures_prodi_id ON fee_structures(prodi_id);
CREATE INDEX idx_fee_structures_academic_year ON fee_structures(academic_year);
CREATE INDEX idx_fee_structures_is_active ON fee_structures(is_active);

-- Trigger for updated_at
CREATE TRIGGER trigger_fee_structures_updated_at
    BEFORE UPDATE ON fee_structures
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
