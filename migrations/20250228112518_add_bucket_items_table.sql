-- +goose Up
-- +goose StatementBegin
CREATE TABLE storage.files (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bucket_id UUID REFERENCES storage.buckets(uuid) ON DELETE CASCADE,
    name varchar NOT NULL,
    path TEXT NOT NULL,
    size BIGINT NOT NULL,
    mime_type TEXT NOT NULL,
    created_by UUID NOT NULL REFERENCES authentication.users(uuid) ON DELETE CASCADE,
    updated_by UUID NOT NULL REFERENCES authentication.users(uuid) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (bucket_id, path)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE storage.files;
-- +goose StatementEnd
