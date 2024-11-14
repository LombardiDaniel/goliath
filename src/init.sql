CREATE TABLE IF NOT EXISTS sessions (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    token VARCHAR(64) UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    -- Index for token lookups
    CONSTRAINT sessions_token_key UNIQUE (token),
    -- Index for cleanup and user lookups
    INDEX sessions_expires_at_idx (expires_at),
    INDEX sessions_user_id_idx (user_id)
);

-- Create function to delete expired sessions
CREATE OR REPLACE FUNCTION delete_expired_sessions()
RETURNS trigger AS $$
BEGIN
    DELETE FROM sessions
    WHERE expires_at < NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger that runs periodically
DROP TRIGGER IF EXISTS cleanup_expired_sessions ON sessions;
CREATE TRIGGER cleanup_expired_sessions
    AFTER INSERT ON sessions
    EXECUTE FUNCTION delete_expired_sessions();