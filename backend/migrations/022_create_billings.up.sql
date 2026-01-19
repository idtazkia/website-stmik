-- Billings table: tracks what candidates owe
CREATE TABLE billings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_candidate UUID NOT NULL REFERENCES candidates(id) ON DELETE CASCADE,
    billing_type VARCHAR(50) NOT NULL, -- registration, tuition, dormitory
    description TEXT,
    amount INTEGER NOT NULL, -- in IDR (smallest unit)
    due_date DATE,
    status VARCHAR(20) NOT NULL DEFAULT 'unpaid', -- unpaid, pending_verification, paid, cancelled
    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_billings_candidate ON billings(id_candidate);
CREATE INDEX idx_billings_status ON billings(status);
CREATE INDEX idx_billings_type ON billings(billing_type);

-- Seed billing types for reference (stored as varchar in billings table)
-- Valid types: registration, tuition, dormitory
COMMENT ON COLUMN billings.billing_type IS 'Valid types: registration, tuition, dormitory';
