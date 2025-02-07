-- +goose Up
-- +goose StatementBegin
CREATE TABLE fluxton.tables (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES fluxton.projects(id) ON DELETE CASCADE,
    created_by UUID NOT NULL REFERENCES authentication.users(id),
    updated_by UUID NOT NULL REFERENCES authentication.users(id),
    name VARCHAR(255) NOT NULL,
    columns JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE fluxton.tables;
-- +goose StatementEnd
