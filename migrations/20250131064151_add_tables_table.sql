-- +goose Up
-- +goose StatementBegin
CREATE TABLE tables (
    id SERIAL PRIMARY KEY,
    project_id INT NOT NULL,
    created_by INT NOT NULL,
    updated_by INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    columns JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE tables ADD CONSTRAINT fk_project_id FOREIGN KEY (project_id) REFERENCES projects(id);
ALTER TABLE tables ADD CONSTRAINT fk_created_by FOREIGN KEY (created_by) REFERENCES users(id);
ALTER TABLE tables ADD CONSTRAINT fk_updated_by FOREIGN KEY (updated_by) REFERENCES users(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tables;
-- +goose StatementEnd
