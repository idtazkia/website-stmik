-- Remove NIM and referral_code fields from candidates table
DROP INDEX IF EXISTS idx_candidates_referral_code;
DROP INDEX IF EXISTS idx_candidates_nim;
ALTER TABLE candidates DROP COLUMN IF EXISTS enrolled_at;
ALTER TABLE candidates DROP COLUMN IF EXISTS referral_code;
ALTER TABLE candidates DROP COLUMN IF EXISTS nim;
