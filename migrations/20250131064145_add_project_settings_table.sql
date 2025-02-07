-- +goose Up
-- +goose StatementBegin
CREATE TABLE fluxton.project_settings (
    id SERIAL PRIMARY KEY,
    project_id UUID NOT NULL REFERENCES fluxton.projects (id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    value TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE fluxton.project_settings;
-- +goose StatementEnd
