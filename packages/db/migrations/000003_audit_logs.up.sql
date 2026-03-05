-- Audit logs table
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    action TEXT NOT NULL,         -- e.g., 'READ_SECRET', 'CREATE_PROJECT'
    target_type TEXT NOT NULL,    -- e.g., 'SECRET', 'PROJECT', 'CLIENT'
    target_id UUID,               -- The ID of the object being acted upon
    metadata JSONB,               -- Extra info like IP, User-Agent, etc.
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index for faster querying of logs
CREATE INDEX idx_audit_logs_target ON audit_logs (target_type, target_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs (created_at DESC);
