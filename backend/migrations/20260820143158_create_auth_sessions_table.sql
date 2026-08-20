-- +goose Up

CREATE TABLE auth_sessions
(
    id         BIGSERIAL PRIMARY KEY,

    user_id    BIGINT      NOT NULL
        REFERENCES users (id)
            ON DELETE CASCADE,

    token_hash CHAR(64)    NOT NULL UNIQUE,

    expires_at TIMESTAMPTZ NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_auth_sessions_user_id
    ON auth_sessions (user_id);

CREATE INDEX idx_auth_sessions_expires_at
    ON auth_sessions (expires_at);

-- +goose Down

DROP TABLE auth_sessions;