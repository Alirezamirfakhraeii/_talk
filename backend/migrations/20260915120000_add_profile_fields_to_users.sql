-- +goose Up
ALTER TABLE users
    ADD COLUMN bio VARCHAR(160),
    ADD COLUMN avatar_path TEXT;

-- +goose Down
ALTER TABLE users
DROP COLUMN avatar_path,
    DROP COLUMN bio;