-- +goose Up
-- +goose StatementBegin
CREATE TABLE fluxton.organization_users (
    id SERIAL PRIMARY KEY,
    organization_id UUID NOT NULL REFERENCES fluxton.organizations (id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES authentication.users (id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE fluxton.organization_users;
-- +goose StatementEnd
