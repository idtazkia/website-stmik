-- Create sessions table for refresh token storage
CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token VARCHAR(500) NOT NULL UNIQUE,
    user_agent VARCHAR(500),
    ip_address VARCHAR(45),
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create index on user_id for user's sessions lookup
CREATE INDEX idx_sessions_user_id ON sessions(user_id);

-- Create index on refresh_token for token validation
CREATE INDEX idx_sessions_refresh_token ON sessions(refresh_token);

-- Create index on expires_at for cleanup job
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);
