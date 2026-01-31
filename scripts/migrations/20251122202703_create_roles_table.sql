-- +goose Up
-- +goose StatementBegin
CREATE TABLE roles
(
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT,
    description TEXT,
    role        TEXT        NOT NULL UNIQUE,
    service_id  BIGINT      NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (service_id, role)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE roles;
-- +goose StatementEnd
