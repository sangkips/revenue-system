CREATE TABLE revenues (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    taxpayer_id UUID NOT NULL REFERENCES taxpayers(id) ON DELETE RESTRICT,
    county_id INTEGER NOT NULL REFERENCES counties(id) ON DELETE RESTRICT,
    amount DECIMAL(15,2) NOT NULL CHECK (amount >= 0),
    revenue_type TEXT NOT NULL CHECK (revenue_type IN ('tax', 'fee', 'fine')),
    transaction_date DATE NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);