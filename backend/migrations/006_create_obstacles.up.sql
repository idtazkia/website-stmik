-- Obstacles that prevent candidate conversion
CREATE TABLE obstacles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    suggested_response TEXT,
    is_active BOOLEAN DEFAULT true,
    display_order INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index
CREATE INDEX idx_obstacles_is_active ON obstacles(is_active);
CREATE INDEX idx_obstacles_display_order ON obstacles(display_order);
