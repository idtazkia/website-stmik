-- Assignment algorithms for candidate-to-consultant assignment
CREATE TABLE assignment_algorithms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL UNIQUE,
    code VARCHAR(30) NOT NULL UNIQUE,
    description TEXT,
    is_active BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Ensure only one algorithm is active at a time
CREATE UNIQUE INDEX idx_assignment_algorithms_single_active
ON assignment_algorithms (is_active)
WHERE is_active = true;

-- Seed default algorithms
INSERT INTO assignment_algorithms (name, code, description, is_active) VALUES
('Round Robin', 'round_robin', 'Assign candidates to consultants in sequential order, cycling through all active consultants equally', true),
('Load Balanced', 'load_balanced', 'Assign to consultant with fewest active candidates (prospecting + committed status)', false),
('Performance Weighted', 'performance_weighted', 'Assign more candidates to consultants with higher enrollment conversion rates', false),
('Activity Based', 'activity_based', 'Assign to consultant with most recent interaction activity, favoring active consultants', false);
