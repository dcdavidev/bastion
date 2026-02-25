-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    salt TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'COLLABORATOR', -- 'ADMIN' or 'COLLABORATOR'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Access control: stores per-user wrapped data keys for each project
CREATE TABLE IF NOT EXISTS user_project_access (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    wrapped_data_key TEXT NOT NULL,
    PRIMARY KEY (user_id, project_id)
);

-- Create an index for faster access checks
CREATE INDEX idx_user_project_access_user ON user_project_access(user_id);
