-- Create payments table
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    county_id INTEGER NOT NULL REFERENCES counties(id) ON DELETE RESTRICT,
    taxpayer_id UUID NOT NULL REFERENCES taxpayers(id) ON DELETE RESTRICT,
    assessment_id UUID REFERENCES assessments(id) ON DELETE SET NULL,
    payment_number TEXT NOT NULL UNIQUE,
    amount DECIMAL(15,2) NOT NULL CHECK (amount > 0),
    payment_method TEXT NOT NULL CHECK (payment_method IN ('mpesa', 'bank_transfer', 'card', 'cheque', 'cash')),
    payment_channel TEXT,  -- e.g., 'mobile_app', 'web'
    external_transaction_id TEXT,
    payer_phone_number TEXT,
    payer_name TEXT,
    payment_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed', 'failed', 'cancelled')),
    collected_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_payments_taxpayer ON payments(taxpayer_id);
CREATE INDEX idx_payments_assessment ON payments(assessment_id);
CREATE INDEX idx_payments_date ON payments(payment_date);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_method ON payments(payment_method);