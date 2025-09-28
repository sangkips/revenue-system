-- Return the created user
-- name: InsertUser :one
INSERT INTO users (id, county_id, username, email, password_hash, first_name, last_name, phone_number, role, employee_id, department, is_active)
VALUES (uuid_generate_v4(), @county_id, @username, @email, @password_hash, @first_name, @last_name, @phone_number, @role, @employee_id, @department, @is_active)
RETURNING id, county_id, username, email, password_hash, first_name, last_name, phone_number, role, employee_id, department, is_active, last_login, created_at, updated_at;

-- name: GetUserByID :one
SELECT id, county_id, username, email, first_name, last_name, phone_number, role, employee_id, department, is_active, last_login, created_at, updated_at
FROM users
WHERE id = @id;

-- name: GetUserByUsername :one
SELECT id, county_id, username, email, password_hash, first_name, last_name, phone_number, role, employee_id, department, is_active, last_login, created_at, updated_at
FROM users
WHERE username = @username;

-- name: ListUsers :many
SELECT id, county_id, username, email, first_name, last_name, phone_number, role, employee_id, department, is_active, last_login, created_at, updated_at
FROM users
WHERE county_id = @county_id
ORDER BY username ASC
LIMIT $1 OFFSET $2;

-- name: UpdateUser :exec
UPDATE users
SET email = @email, first_name = @first_name, last_name = @last_name, phone_number = @phone_number, role = @role, employee_id = @employee_id, department = @department, is_active = @is_active, updated_at = CURRENT_TIMESTAMP
WHERE id = @id;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = @password_hash, updated_at = CURRENT_TIMESTAMP
WHERE id = @id;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = @id;