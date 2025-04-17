-- Create user_sessions table
CREATE TABLE IF NOT EXISTS user_sessions (
    id SERIAL PRIMARY KEY,
    phone_number VARCHAR(20) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    state VARCHAR(50),
    CONSTRAINT unique_phone_number UNIQUE (phone_number)
);

-- Create index on phone_number for faster lookups
CREATE INDEX IF NOT EXISTS idx_user_sessions_phone_number ON user_sessions(phone_number);

-- Add comment to table
COMMENT ON TABLE user_sessions IS 'Stores user session information for WhatsApp users';
