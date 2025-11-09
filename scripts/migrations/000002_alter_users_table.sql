-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN username VARCHAR(32)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
    DROP COLUMN username;
-- +goose StatementEnd
