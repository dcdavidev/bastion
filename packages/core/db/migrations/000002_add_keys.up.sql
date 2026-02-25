-- Table to store global vault configuration (Single User)
CREATE TABLE IF NOT EXISTS vault_config (
    id SERIAL PRIMARY KEY,
    wrapped_master_key TEXT NOT NULL, -- MK encrypted by Admin KEK
    master_key_salt TEXT NOT NULL,    -- Salt used to derive Admin KEK
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Add wrapped_data_key to projects
ALTER TABLE projects ADD COLUMN IF NOT EXISTS wrapped_data_key TEXT;
