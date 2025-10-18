-- Create applications table
CREATE TABLE applications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    taxpayer_id UUID NOT NULL REFERENCES taxpayers(id) ON DELETE RESTRICT,
    type TEXT NOT NULL CHECK (type IN ('single_business_permit', 'building_approval', 'seasonal_parking_ticket', 'health_certificate')),
    notes TEXT,
    status TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'submitted', 'under_review', 'approved', 'rejected')),
    submission_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    approval_date TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create application_documents table
CREATE TABLE application_documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    application_id UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    file_path TEXT NOT NULL, -- e.g., S3 URL or local path
    file_type TEXT NOT NULL CHECK (file_type IN ('pdf', 'jpg', 'png')),
    uploaded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create single_business_permits table
CREATE TABLE single_business_permits (
    application_id UUID PRIMARY KEY REFERENCES applications(id) ON DELETE CASCADE,
    business_name TEXT NOT NULL,
    kra_pin TEXT NOT NULL CHECK (kra_pin ~ '^[A-Za-z][0-9]{9}[A-Za-z]$'),
    business_type TEXT NOT NULL CHECK (business_type IN ('retail_shop', 'hotel', 'wholesale', 'manufacturer')),
    business_location TEXT NOT NULL,
    number_of_employees INTEGER NOT NULL CHECK (number_of_employees >= 0)
);

-- Create building_approvals table
CREATE TABLE building_approvals (
    application_id UUID PRIMARY KEY REFERENCES applications(id) ON DELETE CASCADE,
    project_name TEXT NOT NULL,
    plot_parcel_number TEXT NOT NULL,
    project_type TEXT NOT NULL CHECK (project_type IN ('residential', 'commercial', 'industrial')),
    estimated_project_cost DECIMAL(15,2) NOT NULL CHECK (estimated_project_cost >= 0),
    contact_email TEXT,
    contact_phone TEXT
);

-- Create seasonal_parking_tickets table
CREATE TABLE seasonal_parking_tickets (
    application_id UUID PRIMARY KEY REFERENCES applications(id) ON DELETE CASCADE,
    vehicle_registration_number TEXT NOT NULL CHECK (vehicle_registration_number ~ '^[A-Z0-9]{1,8}$'),
    preferred_parking_zone TEXT NOT NULL,
    duration TEXT NOT NULL CHECK (duration IN ('monthly', 'quarterly', 'annual')),
    contact_email TEXT,
    contact_phone TEXT
);

-- Create health_certificates table
CREATE TABLE health_certificates (
    application_id UUID PRIMARY KEY REFERENCES applications(id) ON DELETE CASCADE,
    applicant_name TEXT NOT NULL,
    business_name TEXT NOT NULL,
    contact_email TEXT,
    contact_phone TEXT
);

-- Create link to assessments
CREATE TABLE application_assessments (
    application_id UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    assessment_id UUID NOT NULL REFERENCES assessments(id) ON DELETE CASCADE,
    PRIMARY KEY (application_id, assessment_id)
);

-- Indexes for performance (unchanged)
CREATE INDEX idx_applications_taxpayer ON applications(taxpayer_id);
CREATE INDEX idx_applications_status ON applications(status);
CREATE INDEX idx_application_documents_application ON application_documents(application_id);