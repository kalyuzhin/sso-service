-- +goose Up
-- +goose StatementBegin
CREATE TABLE refresh_sessions
(
    id                 BIGSERIAL PRIMARY KEY,
    user_id            BIGINT       NOT NULL,
    refresh_token_hash TEXT         NOT NULL,
    user_agent         VARCHAR(200) NOT NULL,
    ip                 VARCHAR(20)  NOT NULL,
    expires_in         TIMESTAMPTZ  NOT NULL,
    created_at         TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE refresh_sessions;
-- +goose StatementEnd
