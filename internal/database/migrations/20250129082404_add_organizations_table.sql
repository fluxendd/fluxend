-- +goose Up
-- +goose StatementBegin
CREATE TABLE fluxend.organizations (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    created_by UUID NOT NULL REFERENCES authentication.users(uuid),
    updated_by UUID NOT NULL REFERENCES authentication.users(uuid),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE fluxend.organizations;
-- +goose StatementEnd
