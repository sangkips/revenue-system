-- name: InsertTaxpayer :one
INSERT INTO taxpayers (
    county_id, user_id, taxpayer_type, national_id, email, phone_number, 
    first_name, last_name, business_name
)
VALUES (
    @county_id, @user_id, @taxpayer_type, @national_id, @email, @phone_number,
    @first_name, @last_name, @business_name
)
RETURNING id, user_id, county_id, taxpayer_type, national_id, email, phone_number, first_name, last_name, business_name, created_at, updated_at;

-- name: GetTaxpayerByID :one
SELECT id, county_id, user_id, taxpayer_type, national_id, email, phone_number,
       first_name, last_name, business_name, created_at, updated_at
FROM taxpayers
WHERE id = @id;

-- name: GetTaxpayerByNationalID :one
SELECT id, county_id, user_id, taxpayer_type, national_id, email, phone_number,
       first_name, last_name, business_name, created_at, updated_at
FROM taxpayers
WHERE national_id = @national_id;

-- name: GetTaxpayerByUserID :one
SELECT id, county_id, user_id, taxpayer_type, national_id, email, phone_number,
       first_name, last_name, business_name, created_at, updated_at
FROM taxpayers
WHERE user_id = @user_id;

-- name: ListTaxpayers :many
SELECT id, county_id, user_id, taxpayer_type, national_id, email, phone_number,
       first_name, last_name, business_name, created_at, updated_at
FROM taxpayers
WHERE county_id = @county_id
ORDER BY created_at ASC
LIMIT $1 OFFSET $2;

-- name: GetFullProfileByUserID :one
SELECT u.id, u.email as u_email, u.first_name as u_first_name, u.last_name as u_last_name, u.phone_number as u_phone,
       u.role, u.created_at as u_created_at,
       t.county_id, t.taxpayer_type, t.national_id, t.email as t_email, t.phone_number as t_phone,
       t.first_name, t.last_name, t.business_name, t.created_at as t_created_at, t.updated_at
FROM users u
JOIN taxpayers t ON u.id = t.user_id
WHERE u.id = @user_id AND u.is_active = true;

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

