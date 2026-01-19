-- Interactions table for logging candidate communications
CREATE TABLE interactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_candidate UUID NOT NULL REFERENCES candidates(id) ON DELETE CASCADE,
    id_consultant UUID NOT NULL REFERENCES users(id),
    channel VARCHAR(50) NOT NULL, -- call, whatsapp, email, campus_visit, home_visit
    id_category UUID REFERENCES interaction_categories(id),
    id_obstacle UUID REFERENCES obstacles(id),
    remarks TEXT NOT NULL,
    next_followup_date DATE,
    next_action TEXT,
    supervisor_suggestion TEXT,
    suggestion_read_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for common queries
CREATE INDEX idx_interactions_id_candidate ON interactions(id_candidate);
CREATE INDEX idx_interactions_id_consultant ON interactions(id_consultant);
CREATE INDEX idx_interactions_created_at ON interactions(created_at DESC);
CREATE INDEX idx_interactions_followup ON interactions(next_followup_date) WHERE next_followup_date IS NOT NULL;

-- Add last_interaction_at to candidates for quick access
ALTER TABLE candidates ADD COLUMN IF NOT EXISTS last_interaction_at TIMESTAMPTZ;

-- Create index for overdue followups query
CREATE INDEX idx_candidates_last_interaction ON candidates(last_interaction_at);
