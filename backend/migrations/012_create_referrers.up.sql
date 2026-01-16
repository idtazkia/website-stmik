-- Referrer Management
-- Tracks referrers (alumni, teachers, partners, staff) for commission tracking

CREATE TABLE referrers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('alumni', 'teacher', 'student', 'partner', 'staff')),
    institution VARCHAR(200),
    phone VARCHAR(20),
    email VARCHAR(100),
    code VARCHAR(20) UNIQUE,
    bank_name VARCHAR(50),
    bank_account VARCHAR(30),
    account_holder VARCHAR(100),
    commission_override BIGINT,
    payout_preference VARCHAR(20) NOT NULL DEFAULT 'per_enrollment' CHECK (payout_preference IN ('monthly', 'per_enrollment')),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for referrer lookup by code
CREATE UNIQUE INDEX idx_referrers_code ON referrers(code) WHERE code IS NOT NULL;

-- Index for referrer lookup by type
CREATE INDEX idx_referrers_type ON referrers(type);

-- Index for active referrers
CREATE INDEX idx_referrers_active ON referrers(is_active) WHERE is_active = true;

-- Index for name search
CREATE INDEX idx_referrers_name ON referrers(name);

COMMENT ON TABLE referrers IS 'Referrers who bring in candidates for commission';
COMMENT ON COLUMN referrers.type IS 'Referrer type: alumni, teacher, student, partner, staff';
COMMENT ON COLUMN referrers.code IS 'Optional referral code for tracking (auto-generated or custom)';
COMMENT ON COLUMN referrers.commission_override IS 'Override default commission amount (null = use reward_config)';
COMMENT ON COLUMN referrers.payout_preference IS 'Payout timing: monthly, per_enrollment';
