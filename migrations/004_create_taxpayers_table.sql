CREATE TABLE taxpayers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    county_id INTEGER NOT NULL REFERENCES counties(id) ON DELETE RESTRICT,
    taxpayer_type TEXT NOT NULL CHECK (taxpayer_type IN ('individual', 'business')),
    national_id TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL,
    phone_number TEXT,
    first_name TEXT, -- Nullable for businesses
    last_name TEXT,  -- Nullable for businesses
    business_name TEXT, -- Nullable for individuals
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);