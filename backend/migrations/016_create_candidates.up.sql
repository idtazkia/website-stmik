-- Candidates table for prospective students
CREATE TABLE candidates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE,
    email_verified BOOLEAN DEFAULT false,
    phone VARCHAR(20) UNIQUE,
    phone_verified BOOLEAN DEFAULT false,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    address TEXT,
    city VARCHAR(100),
    province VARCHAR(100),
    high_school VARCHAR(255),
    graduation_year INT,
    prodi_id UUID REFERENCES prodis(id),
    campaign_id UUID REFERENCES campaigns(id),
    referrer_id UUID REFERENCES referrers(id),
    referred_by_candidate_id UUID REFERENCES candidates(id),
    source_type VARCHAR(50),
    source_detail TEXT,
    assigned_consultant_id UUID REFERENCES users(id),
    status VARCHAR(20) DEFAULT 'registered' CHECK (status IN ('registered', 'prospecting', 'committed', 'enrolled', 'lost')),
    lost_reason_id UUID REFERENCES lost_reasons(id),
    lost_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT email_or_phone_required CHECK (email IS NOT NULL OR phone IS NOT NULL)
);

-- Partial indexes for nullable unique columns
CREATE INDEX idx_candidates_email ON candidates(email) WHERE email IS NOT NULL;
CREATE INDEX idx_candidates_phone ON candidates(phone) WHERE phone IS NOT NULL;
CREATE INDEX idx_candidates_status ON candidates(status);
CREATE INDEX idx_candidates_assigned ON candidates(assigned_consultant_id);
CREATE INDEX idx_candidates_prodi ON candidates(prodi_id);
CREATE INDEX idx_candidates_campaign ON candidates(campaign_id);
CREATE INDEX idx_candidates_referrer ON candidates(referrer_id);

-- Trigger for updated_at
CREATE TRIGGER update_candidates_updated_at
    BEFORE UPDATE ON candidates
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
