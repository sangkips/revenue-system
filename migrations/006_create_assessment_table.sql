-- Create assessments table
CREATE TABLE assessments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    county_id INTEGER NOT NULL REFERENCES counties(id) ON DELETE RESTRICT,
    taxpayer_id UUID NOT NULL REFERENCES taxpayers(id) ON DELETE RESTRICT,
    revenue_id UUID REFERENCES revenues(id) ON DELETE SET NULL,
    assessment_number TEXT NOT NULL UNIQUE,
    assessment_type TEXT NOT NULL,
    financial_year TEXT NOT NULL,  -- e.g., '2025/2026'
    base_amount DECIMAL(15,2) NOT NULL CHECK (base_amount >= 0),
    calculated_amount DECIMAL(15,2) NOT NULL CHECK (calculated_amount >= 0),
    total_amount DECIMAL(15,2) NOT NULL CHECK (total_amount >= 0),
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected', 'paid')),
    due_date DATE NOT NULL,
    assessed_by UUID REFERENCES users(id) ON DELETE SET NULL,
    assessed_date DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create assessment_items table
CREATE TABLE assessment_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    assessment_id UUID NOT NULL REFERENCES assessments(id) ON DELETE CASCADE,
    item_description TEXT NOT NULL,
    quantity DECIMAL(10,2) DEFAULT 1 CHECK (quantity > 0),
    unit_amount DECIMAL(15,2) NOT NULL CHECK (unit_amount >= 0),
    total_amount DECIMAL(15,2) NOT NULL CHECK (total_amount >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_assessments_taxpayer ON assessments(taxpayer_id);
CREATE INDEX idx_assessments_status ON assessments(status);
CREATE INDEX idx_assessments_due_date ON assessments(due_date);
CREATE INDEX idx_assessment_items_assessment ON assessment_items(assessment_id);