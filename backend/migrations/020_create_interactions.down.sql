-- Remove last_interaction_at column from candidates
ALTER TABLE candidates DROP COLUMN IF EXISTS last_interaction_at;

-- Drop interactions table
DROP TABLE IF EXISTS interactions;
