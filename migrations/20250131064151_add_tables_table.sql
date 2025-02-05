-- +goose Up
-- +goose StatementBegin
CREATE TABLE fluxton.tables (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL,
    created_by UUID NOT NULL,
    updated_by UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    columns JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE fluxton.tables ADD CONSTRAINT fk_project_id FOREIGN KEY (project_id) REFERENCES fluxton.projects(id);
ALTER TABLE fluxton.tables ADD CONSTRAINT fk_created_by FOREIGN KEY (created_by) REFERENCES authentication.users(id);
ALTER TABLE fluxton.tables ADD CONSTRAINT fk_updated_by FOREIGN KEY (updated_by) REFERENCES authentication.users(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE fluxton.tables;
-- +goose StatementEnd
