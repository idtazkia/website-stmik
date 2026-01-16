-- Announcements for candidates
CREATE TABLE announcements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    -- Target filters (null means all)
    target_status VARCHAR(20),  -- registered, prospecting, committed, enrolled (null = all statuses)
    target_prodi_id UUID REFERENCES prodis(id),  -- null = all programs
    -- Publishing
    is_published BOOLEAN DEFAULT false,
    published_at TIMESTAMPTZ,
    -- Metadata
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_announcements_published ON announcements(is_published, published_at DESC) WHERE is_published = true;
CREATE INDEX idx_announcements_target_status ON announcements(target_status) WHERE target_status IS NOT NULL;
CREATE INDEX idx_announcements_target_prodi ON announcements(target_prodi_id) WHERE target_prodi_id IS NOT NULL;

-- Track which candidates have read which announcements
CREATE TABLE announcement_reads (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    announcement_id UUID NOT NULL REFERENCES announcements(id) ON DELETE CASCADE,
    candidate_id UUID NOT NULL REFERENCES candidates(id) ON DELETE CASCADE,
    read_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(announcement_id, candidate_id)
);

CREATE INDEX idx_announcement_reads_candidate ON announcement_reads(candidate_id);
