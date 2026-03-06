-- Change value column from JSONB to TEXT in secrets table
ALTER TABLE secrets ALTER COLUMN value TYPE TEXT USING value::TEXT;
