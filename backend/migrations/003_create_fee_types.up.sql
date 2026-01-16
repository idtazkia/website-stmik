-- Fee types table (seeded: registration, tuition, dormitory)
CREATE TABLE fee_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL UNIQUE,
    code VARCHAR(20) NOT NULL UNIQUE,
    is_recurring BOOLEAN DEFAULT false,
    installment_options JSONB DEFAULT '[1]',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index
CREATE INDEX idx_fee_types_code ON fee_types(code);
