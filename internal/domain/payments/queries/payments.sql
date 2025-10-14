-- internal/domains/payments/queries/payments.sql
-- name: InsertPayment :one
INSERT INTO payments (
    county_id, taxpayer_id, assessment_id, payment_number, amount, payment_method,
    payment_channel, external_transaction_id, mpesa_receipt_number, bank_reference,
    cheque_number, payer_phone_number, payer_name, status, collected_by,
    failure_reason, collection_point, gps_coordinates, blockchain_hash, block_number,
    reconciled, reconciliation_date, reconciled_by
)
VALUES (
    @county_id, @taxpayer_id, @assessment_id, @payment_number, @amount, @payment_method,
    @payment_channel, @external_transaction_id, @mpesa_receipt_number, @bank_reference,
    @cheque_number, @payer_phone_number, @payer_name, @status, @collected_by,
    @failure_reason, @collection_point, 
    CASE WHEN @gps_coordinates = '' THEN NULL ELSE @gps_coordinates::point END, 
    @blockchain_hash, @block_number,
    @reconciled, @reconciliation_date, @reconciled_by
)
RETURNING id, county_id, taxpayer_id, assessment_id, payment_number, amount, payment_method,
    payment_channel, external_transaction_id, payer_phone_number, payer_name, payment_date,
    status, collected_by, created_at, updated_at, mpesa_receipt_number, bank_reference,
    cheque_number, failure_reason, collection_point, gps_coordinates, blockchain_hash,
    block_number, reconciled, reconciliation_date, reconciled_by;

-- name: GetPaymentByID :one
SELECT id, county_id, taxpayer_id, assessment_id, payment_number, amount, payment_method,
       payment_channel, external_transaction_id, payer_phone_number, payer_name, payment_date,
       status, collected_by, created_at, updated_at, mpesa_receipt_number, bank_reference,
       cheque_number, failure_reason, collection_point, gps_coordinates, blockchain_hash,
       block_number, reconciled, reconciliation_date, reconciled_by
FROM payments
WHERE id = @id;

-- name: ListPayments :many
SELECT id, county_id, taxpayer_id, assessment_id, payment_number, amount, payment_method,
       payment_channel, external_transaction_id, payer_phone_number, payer_name, payment_date,
       status, collected_by, created_at, updated_at, mpesa_receipt_number, bank_reference,
       cheque_number, failure_reason, collection_point, gps_coordinates, blockchain_hash,
       block_number, reconciled, reconciliation_date, reconciled_by
FROM payments
WHERE county_id = @county_id
ORDER BY payment_date DESC
LIMIT $1 OFFSET $2;

-- name: UpdatePayment :one
UPDATE payments
SET 
    amount = COALESCE(@amount, amount),
    payment_method = COALESCE(@payment_method, payment_method),
    payment_channel = COALESCE(@payment_channel, payment_channel),
    external_transaction_id = COALESCE(@external_transaction_id, external_transaction_id),
    mpesa_receipt_number = COALESCE(@mpesa_receipt_number, mpesa_receipt_number),
    bank_reference = COALESCE(@bank_reference, bank_reference),
    cheque_number = COALESCE(@cheque_number, cheque_number),
    payer_phone_number = COALESCE(@payer_phone_number, payer_phone_number),
    payer_name = COALESCE(@payer_name, payer_name),
    status = COALESCE(@status, status),
    failure_reason = COALESCE(@failure_reason, failure_reason),
    collected_by = COALESCE(@collected_by, collected_by),
    collection_point = COALESCE(@collection_point, collection_point),
    gps_coordinates = CASE 
        WHEN @gps_coordinates = '' THEN gps_coordinates 
        WHEN @gps_coordinates IS NULL THEN gps_coordinates
        ELSE @gps_coordinates::point 
    END,
    blockchain_hash = COALESCE(@blockchain_hash, blockchain_hash),
    block_number = COALESCE(@block_number, block_number),
    reconciled = COALESCE(@reconciled, reconciled),
    reconciliation_date = COALESCE(@reconciliation_date, reconciliation_date),
    reconciled_by = COALESCE(@reconciled_by, reconciled_by),
    updated_at = CURRENT_TIMESTAMP
WHERE id = @id
RETURNING id, county_id, taxpayer_id, assessment_id, payment_number, amount, payment_method,
    payment_channel, external_transaction_id, payer_phone_number, payer_name, payment_date,
    status, collected_by, created_at, updated_at, mpesa_receipt_number, bank_reference,
    cheque_number, failure_reason, collection_point, gps_coordinates, blockchain_hash,
    block_number, reconciled, reconciliation_date, reconciled_by;

-- name: ListPaymentsByRevenueID :many
SELECT id, county_id, taxpayer_id, assessment_id, payment_number, amount, payment_method,
    payment_channel, external_transaction_id, payer_phone_number, payer_name, payment_date,
    status, collected_by, created_at, updated_at, mpesa_receipt_number, bank_reference,
    cheque_number, failure_reason, collection_point, gps_coordinates, blockchain_hash,
    block_number, reconciled, reconciliation_date, reconciled_by
FROM payments
WHERE assessment_id = @assessment_id
ORDER BY payment_date DESC;

-- name: DeletePayment :exec
DELETE FROM payments WHERE id = @id;

-- Payment Allocations Queries
-- name: InsertPaymentAllocation :one
INSERT INTO payment_allocations (
    payment_id, assessment_id, allocated_amount, allocation_type
)
VALUES (
    @payment_id, @assessment_id, @allocated_amount, @allocation_type
)
RETURNING id, payment_id, assessment_id, allocated_amount, allocation_type, created_at;

-- name: ListPaymentAllocations :many
SELECT id, payment_id, assessment_id, allocated_amount, allocation_type, created_at
FROM payment_allocations
WHERE payment_id = @payment_id
ORDER BY created_at ASC;

-- name: DeletePaymentAllocation :exec
DELETE FROM payment_allocations WHERE id = @id;

-- Receipts Queries
-- name: InsertReceipt :exec
INSERT INTO receipts (
    payment_id, receipt_number, receipt_type, pdf_file_path, pdf_file_size,
    pdf_generated, sms_sent, sms_sent_at, email_sent, email_sent_at,
    blockchain_hash, block_number, blockchain_verified, qr_code_data
)
VALUES (
    @payment_id, @receipt_number, @receipt_type, @pdf_file_path, @pdf_file_size,
    @pdf_generated, @sms_sent, @sms_sent_at, @email_sent, @email_sent_at,
    @blockchain_hash, @block_number, @blockchain_verified, @qr_code_data
);

-- name: GetReceiptByID :one
SELECT id, payment_id, receipt_number, receipt_type, pdf_file_path, pdf_file_size,
       pdf_generated, sms_sent, sms_sent_at, email_sent, email_sent_at,
       blockchain_hash, block_number, blockchain_verified, qr_code_data, created_at
FROM receipts
WHERE id = @id;

-- name: ListReceiptsByPayment :many
SELECT id, payment_id, receipt_number, receipt_type, pdf_file_path, pdf_file_size,
       pdf_generated, sms_sent, sms_sent_at, email_sent, email_sent_at,
       blockchain_hash, block_number, blockchain_verified, qr_code_data, created_at
FROM receipts
WHERE payment_id = @payment_id
ORDER BY created_at ASC;

-- name: UpdateReceipt :exec
UPDATE receipts
SET 
    receipt_type = COALESCE(@receipt_type, receipt_type),
    pdf_file_path = COALESCE(@pdf_file_path, pdf_file_path),
    pdf_file_size = COALESCE(@pdf_file_size, pdf_file_size),
    pdf_generated = COALESCE(@pdf_generated, pdf_generated),
    sms_sent = COALESCE(@sms_sent, sms_sent),
    sms_sent_at = COALESCE(@sms_sent_at, sms_sent_at),
    email_sent = COALESCE(@email_sent, email_sent),
    email_sent_at = COALESCE(@email_sent_at, email_sent_at),
    blockchain_verified = COALESCE(@blockchain_verified, blockchain_verified)
WHERE id = @id;

-- name: DeleteReceipt :exec
DELETE FROM receipts WHERE id = @id;