-- Add email to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS email TEXT UNIQUE;
