-- +goose Up
-- +goose StatementBegin
CREATE TABLE storage.backups (
     uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
     project_uuid UUID NOT NULL REFERENCES fluxend.projects(uuid) ON DELETE CASCADE,
     status VARCHAR NOT NULL,
     error TEXT,
     started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     completed_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS backups;
-- +goose StatementEnd
