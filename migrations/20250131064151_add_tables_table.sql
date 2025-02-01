-- +goose Up
-- +goose StatementBegin
CREATE TABLE tables (
    id SERIAL PRIMARY KEY,
    project_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    columns JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE tables ADD CONSTRAINT fk_project_id FOREIGN KEY (project_id) REFERENCES projects(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tables;
-- +goose StatementEnd
