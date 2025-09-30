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

-- name: ListAllUsers :many
SELECT id, county_id, username, email, first_name, last_name, phone_number, role, employee_id, department, is_active, last_login, created_at, updated_at
FROM users
ORDER BY username ASC
LIMIT $1 OFFSET $2;

-- name: UpdateUser :exec
UPDATE users
SET
  email = CASE WHEN @update_email::boolean THEN @email ELSE email END,
  first_name = CASE WHEN @update_first_name::boolean THEN @first_name ELSE first_name END,
  last_name = CASE WHEN @update_last_name::boolean THEN @last_name ELSE last_name END,
  phone_number = CASE WHEN @update_phone_number::boolean THEN @phone_number ELSE phone_number END,
  role = CASE WHEN @update_role::boolean THEN @role ELSE role END,
  employee_id = CASE WHEN @update_employee_id::boolean THEN @employee_id ELSE employee_id END,
  department = CASE WHEN @update_department::boolean THEN @department ELSE department END,
  is_active = CASE WHEN @update_is_active::boolean THEN @is_active ELSE is_active END,
  updated_at = CURRENT_TIMESTAMP
WHERE id = @id;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = @password_hash, updated_at = CURRENT_TIMESTAMP
WHERE id = @id;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = @id;
