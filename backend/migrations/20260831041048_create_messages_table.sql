-- +goose Up

CREATE TABLE messages
(
    id              BIGSERIAL PRIMARY KEY,

    conversation_id BIGINT      NOT NULL
        REFERENCES direct_conversations (id)
            ON DELETE CASCADE,

    sender_id       BIGINT      NOT NULL
        REFERENCES users (id)
            ON DELETE CASCADE,

    content         TEXT        NOT NULL,

    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_messages_conversation_id
    ON messages (conversation_id);

CREATE INDEX idx_messages_sender_id
    ON messages (sender_id);

CREATE INDEX idx_messages_conversation_created_at
    ON messages (conversation_id, created_at);

-- +goose Down

DROP TABLE messages;