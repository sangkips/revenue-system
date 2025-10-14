-- Alter existing payments table to add missing fields
ALTER TABLE payments 
ADD COLUMN IF NOT EXISTS mpesa_receipt_number VARCHAR(50),
ADD COLUMN IF NOT EXISTS bank_reference VARCHAR(100),
ADD COLUMN IF NOT EXISTS cheque_number VARCHAR(50),
ADD COLUMN IF NOT EXISTS failure_reason TEXT,
ADD COLUMN IF NOT EXISTS collection_point VARCHAR(100),
ADD COLUMN IF NOT EXISTS gps_coordinates POINT,
ADD COLUMN IF NOT EXISTS blockchain_hash VARCHAR(64),
ADD COLUMN IF NOT EXISTS block_number BIGINT,
ADD COLUMN IF NOT EXISTS reconciled BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS reconciliation_date TIMESTAMP,
ADD COLUMN IF NOT EXISTS reconciled_by UUID REFERENCES users(id);

-- Update CHECK constraint for status (add 'refunded')
ALTER TABLE payments DROP CONSTRAINT IF EXISTS payments_status_check;
ALTER TABLE payments ADD CONSTRAINT payments_status_check 
CHECK (status IN ('pending', 'processing', 'completed', 'failed', 'cancelled', 'refunded'));

-- Create payment_allocations table
CREATE TABLE IF NOT EXISTS payment_allocations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    payment_id UUID NOT NULL REFERENCES payments(id) ON DELETE CASCADE,
    assessment_id UUID NOT NULL REFERENCES assessments(id) ON DELETE RESTRICT,
    allocated_amount DECIMAL(15,2) NOT NULL CHECK (allocated_amount > 0),
    allocation_type VARCHAR(20) CHECK (allocation_type IN ('principal', 'penalty', 'interest')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create receipts table
CREATE TABLE IF NOT EXISTS receipts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    payment_id UUID NOT NULL REFERENCES payments(id) ON DELETE CASCADE,
    receipt_number VARCHAR(50) UNIQUE NOT NULL,
    receipt_type VARCHAR(20) CHECK (receipt_type IN ('payment', 'provisional', 'official')),
    
    -- Receipt content
    pdf_file_path VARCHAR(500),
    pdf_file_size INTEGER,
    pdf_generated BOOLEAN DEFAULT FALSE,
    
    -- Delivery tracking
    sms_sent BOOLEAN DEFAULT FALSE,
    sms_sent_at TIMESTAMP,
    email_sent BOOLEAN DEFAULT FALSE,
    email_sent_at TIMESTAMP,
    
    -- Blockchain details
    blockchain_hash VARCHAR(64) UNIQUE NOT NULL,
    block_number BIGINT,
    blockchain_verified BOOLEAN DEFAULT FALSE,
    
    -- QR code for verification
    qr_code_data TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for new tables
CREATE INDEX IF NOT EXISTS idx_payment_allocations_payment ON payment_allocations(payment_id);
CREATE INDEX IF NOT EXISTS idx_payment_allocations_assessment ON payment_allocations(assessment_id);
CREATE INDEX IF NOT EXISTS idx_receipts_payment ON receipts(payment_id);
CREATE INDEX IF NOT EXISTS idx_receipts_blockchain_hash ON receipts(blockchain_hash);