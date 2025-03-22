-- +goose Up
-- +goose StatementBegin
DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'oauth_provider') THEN
CREATE TYPE oauth_provider AS ENUM ('github', 'google');
END IF;
END $$;

CREATE TABLE oauth_connections
(
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    provider oauth_provider NOT NULL,
    provider_id TEXT NOT NULL,
    access_token TEXT,
    refresh_token TEXT,
    expires_at TIMESTAMP,
    PRIMARY KEY (user_id, provider),
    UNIQUE (provider, provider_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS oauth_connections;

DO $$ BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'oauth_provider') THEN
DROP TYPE oauth_provider;
END IF;
END $$;
-- +goose StatementEnd
