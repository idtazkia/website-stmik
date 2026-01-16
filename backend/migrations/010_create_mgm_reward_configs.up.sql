-- Member-Get-Member (MGM) Reward Configuration
-- Rewards for enrolled students who refer new candidates

CREATE TABLE mgm_reward_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    academic_year VARCHAR(9) NOT NULL,
    reward_type VARCHAR(20) NOT NULL CHECK (reward_type IN ('cash', 'tuition_discount', 'merchandise')),
    referrer_amount BIGINT NOT NULL,
    referee_amount BIGINT,
    trigger_event VARCHAR(20) NOT NULL CHECK (trigger_event IN ('commitment', 'enrollment')),
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(academic_year, reward_type, trigger_event)
);

-- Index for active MGM config lookup by academic year
CREATE INDEX idx_mgm_reward_configs_year_active ON mgm_reward_configs(academic_year, is_active) WHERE is_active = true;

COMMENT ON TABLE mgm_reward_configs IS 'Member-Get-Member reward configuration by academic year';
COMMENT ON COLUMN mgm_reward_configs.academic_year IS 'Academic year in format YYYY/YYYY (e.g., 2025/2026)';
COMMENT ON COLUMN mgm_reward_configs.referrer_amount IS 'Reward amount for the referring enrolled student';
COMMENT ON COLUMN mgm_reward_configs.referee_amount IS 'Discount/reward for the new candidate (nullable)';
COMMENT ON COLUMN mgm_reward_configs.trigger_event IS 'When reward is triggered: commitment, enrollment';
