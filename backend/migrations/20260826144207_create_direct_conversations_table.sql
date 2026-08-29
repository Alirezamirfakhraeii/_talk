-- +goose Up

CREATE TABLE direct_conversations
(
    id          BIGSERIAL PRIMARY KEY,

    user_one_id BIGINT      NOT NULL
        REFERENCES users (id)
            ON DELETE CASCADE,

    user_two_id BIGINT      NOT NULL
        REFERENCES users (id)
            ON DELETE CASCADE,

    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT direct_conversations_different_users
        CHECK (user_one_id <> user_two_id),

    CONSTRAINT direct_conversations_unique_pair
        UNIQUE (user_one_id, user_two_id)
);

CREATE INDEX idx_direct_conversations_user_one_id
    ON direct_conversations (user_one_id);

CREATE INDEX idx_direct_conversations_user_two_id
    ON direct_conversations (user_two_id);

-- +goose Down

DROP TABLE direct_conversations;