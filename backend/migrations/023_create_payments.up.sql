-- Payments table: tracks payment proofs uploaded by candidates
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    billing_id UUID NOT NULL REFERENCES billings(id) ON DELETE CASCADE,
    amount INTEGER NOT NULL, -- amount paid in IDR
    transfer_date DATE NOT NULL,
    proof_file_path VARCHAR(500) NOT NULL,
    proof_file_name VARCHAR(255) NOT NULL,
    proof_file_size INTEGER NOT NULL,
    proof_mime_type VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, approved, rejected
    rejection_reason TEXT,
    reviewed_by UUID REFERENCES users(id),
    reviewed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payments_billing ON payments(billing_id);
CREATE INDEX idx_payments_status ON payments(status);
