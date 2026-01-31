-- +goose Up
-- +goose StatementBegin
CREATE TABLE services
(
    id   BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE services;
-- +goose StatementEnd
