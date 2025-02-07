-- +goose Up
-- +goose StatementBegin
CREATE TABLE fluxton.projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES fluxton.organizations (id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    db_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_projects_name_organization_id ON fluxton.projects (name, organization_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE fluxton.projects;
-- +goose StatementEnd
