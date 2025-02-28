-- +goose Up
-- +goose StatementBegin
CREATE TABLE storage.buckets (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name varchar NOT NULL,
    is_public BOOLEAN DEFAULT FALSE,
    project_uuid UUID NOT NULL REFERENCES fluxton.projects(uuid) ON DELETE CASCADE,
    created_by UUID NOT NULL REFERENCES authentication.users(uuid) ON DELETE CASCADE,
    updated_by UUID NOT NULL REFERENCES authentication.users(uuid) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    UNIQUE (project_uuid, name)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE storage.buckets;
-- +goose StatementEnd
