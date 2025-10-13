-- internal/domains/payments/queries/payments.sql
-- name: InsertPayment :one
INSERT INTO payments (
    county_id, taxpayer_id, assessment_id, payment_number, amount, payment_method,
    payment_channel, external_transaction_id, payer_phone_number, payer_name,
    status, collected_by
)
VALUES (
    @county_id, @taxpayer_id, @assessment_id, @payment_number, @amount, @payment_method,
    @payment_channel, @external_transaction_id, @payer_phone_number, @payer_name,
    @status, @collected_by
)
RETURNING id, county_id, taxpayer_id, assessment_id, payment_number, amount, payment_method,
    payment_channel, external_transaction_id, payer_phone_number, payer_name,
    payment_date, status, collected_by, created_at, updated_at;

-- name: GetPaymentByID :one
SELECT id, county_id, taxpayer_id, assessment_id, payment_number, amount, payment_method,
       payment_channel, external_transaction_id, payer_phone_number, payer_name,
       payment_date, status, collected_by, created_at, updated_at
FROM payments
WHERE id = @id;

-- name: ListPayments :many
SELECT id, county_id, taxpayer_id, assessment_id, payment_number, amount, payment_method,
       payment_channel, external_transaction_id, payer_phone_number, payer_name,
       payment_date, status, collected_by, created_at, updated_at
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
    payer_phone_number = COALESCE(@payer_phone_number, payer_phone_number),
    payer_name = COALESCE(@payer_name, payer_name),
    status = COALESCE(@status, status),
    collected_by = COALESCE(@collected_by, collected_by),
    updated_at = CURRENT_TIMESTAMP
WHERE id = @id
RETURNING id, county_id, taxpayer_id, assessment_id, payment_number, amount, payment_method,
    payment_channel, external_transaction_id, payer_phone_number, payer_name,
    payment_date, status, collected_by, created_at, updated_at;

-- name: ListPaymentsByRevenueID :many
SELECT id, county_id, taxpayer_id, assessment_id, payment_number, amount, payment_method,
       payment_channel, external_transaction_id, payer_phone_number, payer_name,
       payment_date, status, collected_by, created_at, updated_at
FROM payments
WHERE assessment_id = @assessment_id
ORDER BY payment_date DESC;

-- name: DeletePayment :exec
DELETE FROM payments WHERE id = @id;