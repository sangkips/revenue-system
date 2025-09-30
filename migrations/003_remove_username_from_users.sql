-- Remove username column from users table
-- This migration removes the username column and its associated index

-- Drop the username index
DROP INDEX IF EXISTS idx_users_username;

-- Drop the username column
ALTER TABLE users DROP COLUMN IF EXISTS username;
