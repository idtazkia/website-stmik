-- Interaction categories for candidate follow-up
CREATE TABLE interaction_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    sentiment VARCHAR(20) NOT NULL CHECK (sentiment IN ('positive', 'neutral', 'negative')),
    is_active BOOLEAN DEFAULT true,
    display_order INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index
CREATE INDEX idx_interaction_categories_sentiment ON interaction_categories(sentiment);
CREATE INDEX idx_interaction_categories_is_active ON interaction_categories(is_active);
CREATE INDEX idx_interaction_categories_display_order ON interaction_categories(display_order);
