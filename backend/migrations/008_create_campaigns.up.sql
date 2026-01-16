-- Campaign Management
-- Campaigns track marketing initiatives for lead generation

CREATE TABLE campaigns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) UNIQUE,
    type VARCHAR(20) NOT NULL CHECK (type IN ('promo', 'event', 'ads', 'organic')),
    channel VARCHAR(50),
    description TEXT,
    start_date DATE,
    end_date DATE,
    registration_fee_override BIGINT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for active campaigns lookup
CREATE INDEX idx_campaigns_active ON campaigns(is_active) WHERE is_active = true;

-- Index for date-based queries
CREATE INDEX idx_campaigns_dates ON campaigns(start_date, end_date);

COMMENT ON TABLE campaigns IS 'Marketing campaigns for tracking lead sources';
COMMENT ON COLUMN campaigns.type IS 'Campaign type: promo, event, ads, organic';
COMMENT ON COLUMN campaigns.channel IS 'Marketing channel: instagram, google, facebook, tiktok, expo, school_visit, etc.';
COMMENT ON COLUMN campaigns.registration_fee_override IS 'Override registration fee (null = use default)';
