-- +goose Up
-- +goose StatementBegin
CREATE TYPE project_status AS ENUM ('active', 'inactive', 'error', 'frozen');

CREATE TABLE fluxton.projects (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_uuid UUID NOT NULL REFERENCES fluxton.organizations (uuid) ON DELETE CASCADE,
    created_by UUID NOT NULL REFERENCES authentication.users (uuid) ON DELETE CASCADE,
    updated_by UUID NOT NULL REFERENCES authentication.users (uuid) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    status project_status NOT NULL DEFAULT 'inactive',
    description TEXT,
    db_name VARCHAR(255) NOT NULL,
    db_port INT NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_projects_name_organization_uuid ON fluxton.projects (name, organization_uuid);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS storage.backups;
DROP TABLE fluxton.projects;
DROP type project_status;
-- +goose StatementEnd
