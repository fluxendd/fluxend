-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_settings (
    id SERIAL PRIMARY KEY,
    project_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    value TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE project_settings ADD CONSTRAINT fk_project_id FOREIGN KEY (project_id) REFERENCES projects(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE project_settings;
-- +goose StatementEnd
