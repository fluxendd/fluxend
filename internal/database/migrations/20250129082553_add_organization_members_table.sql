-- +goose Up
-- +goose StatementBegin
CREATE TABLE fluxend.organization_members (
    id SERIAL PRIMARY KEY,
    organization_uuid UUID NOT NULL REFERENCES fluxend.organizations (uuid) ON DELETE CASCADE,
    user_uuid UUID NOT NULL REFERENCES authentication.users (uuid),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE fluxend.organization_members;
-- +goose StatementEnd
