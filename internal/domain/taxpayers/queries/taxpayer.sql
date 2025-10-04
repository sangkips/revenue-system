-- name: InsertTaxpayer :one
INSERT INTO taxpayers (
    county_id, taxpayer_type, national_id, email, phone_number, 
    first_name, last_name, business_name
)
VALUES (
    @county_id, @taxpayer_type, @national_id, @email, @phone_number,
    @first_name, @last_name, @business_name
)
RETURNING id, county_id, taxpayer_type, national_id, email, phone_number, first_name, last_name, business_name, created_at, updated_at;

-- name: GetTaxpayerByID :one
SELECT id, county_id, taxpayer_type, national_id, email, phone_number,
       first_name, last_name, business_name, created_at, updated_at
FROM taxpayers
WHERE id = @id;

-- name: GetTaxpayerByNationalID :one
SELECT id, county_id, taxpayer_type, national_id, email, phone_number,
       first_name, last_name, business_name, created_at, updated_at
FROM taxpayers
WHERE national_id = @national_id;

-- name: ListTaxpayers :many
SELECT id, county_id, taxpayer_type, national_id, email, phone_number,
       first_name, last_name, business_name, created_at, updated_at
FROM taxpayers
WHERE county_id = @county_id
ORDER BY created_at ASC
LIMIT $1 OFFSET $2;

-- name: UpdateTaxpayer :one
UPDATE taxpayers
SET
    email = @email,
    phone_number = @phone_number,
    first_name = @first_name,
    last_name = @last_name,
    business_name = @business_name,
    updated_at = CURRENT_TIMESTAMP
WHERE id = @id
RETURNING id, county_id, taxpayer_type, national_id, email, phone_number, first_name, last_name, business_name, created_at, updated_at;

-- name: DeleteTaxpayer :exec
DELETE FROM taxpayers WHERE id = @id;
