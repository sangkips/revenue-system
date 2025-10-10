-- name: InsertRevenue :one
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

-- name: UpdateRevenue :one
UPDATE revenues
SET
    amount = COALESCE(sqlc.narg(amount), amount),
    revenue_type = COALESCE(sqlc.narg(revenue_type), revenue_type),
    transaction_date = COALESCE(sqlc.narg(transaction_date), transaction_date),
    description = COALESCE(sqlc.narg(description), description),
    updated_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg(id)
RETURNING id, taxpayer_id, county_id, amount, revenue_type, transaction_date, description, created_at, updated_at;

-- name: ListRevenuesByTaxpayerID :many
SELECT id, taxpayer_id, county_id, amount, revenue_type, transaction_date, description,
       created_at, updated_at
FROM revenues
WHERE taxpayer_id = @taxpayer_id
ORDER BY transaction_date DESC
LIMIT $1 OFFSET $2;

-- name: DeleteRevenue :exec
DELETE FROM revenues WHERE id = @id;
