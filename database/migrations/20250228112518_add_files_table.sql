-- +goose Up
-- +goose StatementBegin
CREATE TABLE storage.files (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    container_uuid UUID NOT NULL REFERENCES storage.containers(uuid) ON DELETE CASCADE,
    full_file_name TEXT NOT NULL,
    size BIGINT NOT NULL,
    mime_type TEXT NOT NULL,
    created_by UUID NOT NULL REFERENCES authentication.users(uuid) ON DELETE CASCADE,
    updated_by UUID NOT NULL REFERENCES authentication.users(uuid) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (container_uuid, full_file_name)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE storage.files;
-- +goose StatementEnd
