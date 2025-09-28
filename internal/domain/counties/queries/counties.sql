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

-- name: UpdateCounty :exec
UPDATE counties
SET name = @name, treasury_account = @treasury_account, updated_at = CURRENT_TIMESTAMP
WHERE id = @id;

-- name: DeleteCounty :exec
DELETE FROM counties WHERE id = @id;
