-- +goose Up
-- +goose StatementBegin
CREATE TABLE projects (
    id SERIAL PRIMARY KEY,
    organization_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    db_name uuid NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE projects ADD CONSTRAINT fk_organization_id FOREIGN KEY (organization_id) REFERENCES organizations(id);
CREATE UNIQUE INDEX idx_projects_name_organization_id ON projects (name, organization_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE projects;
-- +goose StatementEnd
