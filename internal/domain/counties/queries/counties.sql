-- name: InsertCounty :one
INSERT INTO counties (name, code, treasury_account)
VALUES (@name, @code, @treasury_account)
RETURNING id, name, code, treasury_account, created_at, updated_at;

-- name: GetCountyByID :one
SELECT id, name, code, treasury_account, created_at, updated_at
FROM counties
WHERE id = @id;

-- name: ListCounties :many
SELECT id, name, code, treasury_account, created_at, updated_at
FROM counties
ORDER BY name ASC
LIMIT $1 OFFSET $2;

-- name: UpdateCounty :one
UPDATE counties
SET
    name = CASE WHEN @update_name::bool THEN @name ELSE name END,
    treasury_account = CASE WHEN @update_treasury_account::bool THEN @treasury_account ELSE treasury_account END,
    updated_at = CURRENT_TIMESTAMP
WHERE id = @id
RETURNING id, name, code, treasury_account, created_at, updated_at;

-- name: DeleteCounty :exec
DELETE FROM counties WHERE id = @id;
