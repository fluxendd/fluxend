-- +goose Up
-- +goose StatementBegin
CREATE TABLE storage.files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bucket_id UUID REFERENCES storage.buckets(uuid) ON DELETE CASCADE,
    name varchar NOT NULL,
    path TEXT NOT NULL,
    size BIGINT NOT NULL,
    mime_type TEXT NOT NULL,
    uploaded_by UUID NOT NULL REFERENCES authentication.users(uuid) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE storage.files;
-- +goose StatementEnd
