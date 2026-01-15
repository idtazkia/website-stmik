-- Create applications table for admission submissions
CREATE TABLE IF NOT EXISTS applications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    application_number VARCHAR(50) NOT NULL UNIQUE,
    program VARCHAR(50) NOT NULL CHECK (program IN ('SI', 'TI')),
    academic_year VARCHAR(9) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'submitted', 'reviewing', 'accepted', 'rejected', 'enrolled')),

    -- Personal data
    birth_place VARCHAR(100),
    birth_date DATE,
    gender VARCHAR(10) CHECK (gender IN ('male', 'female')),
    address TEXT,
    city VARCHAR(100),
    province VARCHAR(100),
    postal_code VARCHAR(10),

    -- Education background
    school_name VARCHAR(255),
    school_address TEXT,
    graduation_year INTEGER,

    -- Documents
    photo_url VARCHAR(500),
    id_card_url VARCHAR(500),
    diploma_url VARCHAR(500),
    transcript_url VARCHAR(500),

    -- Timestamps
    submitted_at TIMESTAMPTZ,
    reviewed_at TIMESTAMPTZ,
    reviewed_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create index on user_id for user's applications lookup
CREATE INDEX idx_applications_user_id ON applications(user_id);

-- Create index on status for filtering
CREATE INDEX idx_applications_status ON applications(status);

-- Create index on academic_year for reporting
CREATE INDEX idx_applications_academic_year ON applications(academic_year);

-- Create index on program for filtering
CREATE INDEX idx_applications_program ON applications(program);
