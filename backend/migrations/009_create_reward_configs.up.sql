-- Reward Configuration for External Referrers
-- Defines default rewards by referrer type (alumni, teacher, student, partner, staff)

CREATE TABLE reward_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    referrer_type VARCHAR(20) NOT NULL CHECK (referrer_type IN ('alumni', 'teacher', 'student', 'partner', 'staff')),
    reward_type VARCHAR(20) NOT NULL CHECK (reward_type IN ('cash', 'tuition_discount', 'merchandise')),
    amount BIGINT NOT NULL,
    is_percentage BOOLEAN NOT NULL DEFAULT false,
    trigger_event VARCHAR(20) NOT NULL CHECK (trigger_event IN ('registration', 'commitment', 'enrollment')),
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(referrer_type, reward_type, trigger_event)
);

-- Index for active reward config lookup by type
CREATE INDEX idx_reward_configs_type_active ON reward_configs(referrer_type, is_active) WHERE is_active = true;

COMMENT ON TABLE reward_configs IS 'Default reward configuration for external referrers';
COMMENT ON COLUMN reward_configs.referrer_type IS 'Type of referrer: alumni, teacher, student, partner, staff';
COMMENT ON COLUMN reward_configs.reward_type IS 'Type of reward: cash, tuition_discount, merchandise';
COMMENT ON COLUMN reward_configs.amount IS 'Reward amount in IDR (or percentage if is_percentage=true)';
COMMENT ON COLUMN reward_configs.is_percentage IS 'If true, amount is percentage (e.g., 10 = 10%)';
COMMENT ON COLUMN reward_configs.trigger_event IS 'When reward is triggered: registration, commitment, enrollment';
