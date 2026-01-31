-- +goose Up
-- +goose StatementBegin
CREATE TABLE users_roles
(
    id         BIGSERIAL PRIMARY KEY,
    role_id    BIGINT      NOT NULL,
    user_id    BIGINT      NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX ON users_roles(role_id, user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users_roles;
-- +goose StatementEnd
