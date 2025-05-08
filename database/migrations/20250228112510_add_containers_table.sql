-- +goose Up
-- +goose StatementBegin
CREATE TABLE storage.containers (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_uuid UUID NOT NULL REFERENCES fluxton.projects(uuid) ON DELETE CASCADE,
    created_by UUID NOT NULL REFERENCES authentication.users(uuid) ON DELETE CASCADE,
    updated_by UUID NOT NULL REFERENCES authentication.users(uuid) ON DELETE CASCADE,
    name varchar NOT NULL,
    name_key varchar NOT NULL,
    description TEXT,
    is_public BOOLEAN DEFAULT FALSE,
    url varchar,
    total_files BIGINT DEFAULT 0,
    max_file_size BIGINT DEFAULT 2048,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    UNIQUE (project_uuid, name)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE storage.containers;
-- +goose StatementEnd
