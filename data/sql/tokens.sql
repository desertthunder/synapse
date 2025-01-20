-- Token Table
-- Used to authenticate with Discord & BlueSky APIs
CREATE TABLE IF NOT EXISTS tokens (
    id SERIAL PRIMARY KEY,
    token TEXT NOT NULL,
    -- access, refresh
    type VARCHAR(255) NOT NULL,
    -- discord, bluesky
    api VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
