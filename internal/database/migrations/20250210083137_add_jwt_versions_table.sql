-- +goose Up
-- +goose StatementBegin
CREATE TABLE authentication.jwt_versions (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE REFERENCES authentication.users(uuid) ON DELETE CASCADE,
    version INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE authentication.jwt_versions;
-- +goose StatementEnd
