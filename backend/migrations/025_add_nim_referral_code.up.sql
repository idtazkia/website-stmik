-- Add NIM and referral_code fields to candidates table for enrollment and MGM
ALTER TABLE candidates ADD COLUMN nim VARCHAR(20) UNIQUE;
ALTER TABLE candidates ADD COLUMN referral_code VARCHAR(20) UNIQUE;
ALTER TABLE candidates ADD COLUMN enrolled_at TIMESTAMPTZ;

-- Index for referral code lookups
CREATE INDEX idx_candidates_referral_code ON candidates(referral_code) WHERE referral_code IS NOT NULL;

-- Index for NIM lookups
CREATE INDEX idx_candidates_nim ON candidates(nim) WHERE nim IS NOT NULL;
