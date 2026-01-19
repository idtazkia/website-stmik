-- Documents: Store candidate uploaded documents
CREATE TABLE documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_candidate UUID NOT NULL REFERENCES candidates(id) ON DELETE CASCADE,
    id_document_type UUID NOT NULL REFERENCES document_types(id),
    file_name TEXT NOT NULL,
    file_path TEXT NOT NULL,
    file_size INT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    rejection_reason TEXT,
    id_reviewed_by UUID REFERENCES users(id),
    reviewed_at TIMESTAMPTZ,
    is_deferred BOOLEAN DEFAULT false,
    deferred_reason TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT documents_status_check CHECK (status IN ('pending', 'approved', 'rejected'))
);

-- Indexes
CREATE INDEX idx_documents_candidate ON documents(id_candidate);
CREATE INDEX idx_documents_type ON documents(id_document_type);
CREATE INDEX idx_documents_status ON documents(status);

-- Unique constraint: one document per type per candidate (latest upload)
CREATE UNIQUE INDEX idx_documents_candidate_type ON documents(id_candidate, id_document_type);
