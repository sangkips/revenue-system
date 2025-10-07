-- name: InsertRevenue :exec
INSERT INTO revenues (
    taxpayer_id, county_id, amount, revenue_type, transaction_date, description
)
VALUES (
    @taxpayer_id, @county_id, @amount, @revenue_type, @transaction_date, @description
)
RETURNING id, taxpayer_id, county_id, amount, revenue_type, transaction_date, description, created_at, updated_at;

-- name: GetRevenueByID :one
SELECT id, taxpayer_id, county_id, amount, revenue_type, transaction_date, description,
       created_at, updated_at
FROM revenues
WHERE id = @id;

-- name: ListRevenues :many
SELECT id, taxpayer_id, county_id, amount, revenue_type, transaction_date, description,
       created_at, updated_at
FROM revenues
WHERE county_id = @county_id
ORDER BY transaction_date DESC
LIMIT $1 OFFSET $2;

-- name: UpdateRevenue :exec
UPDATE revenues
SET 
    amount = @amount,
    revenue_type = @revenue_type,
    transaction_date = @transaction_date,
    description = @description,
    updated_at = CURRENT_TIMESTAMP
WHERE id = @id
RETURNING id, taxpayer_id, county_id, amount, revenue_type, transaction_date, description, created_at, updated_at;

-- name: DeleteRevenue :exec
DELETE FROM revenues WHERE id = @id;