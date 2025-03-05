-- +goose Up
-- +goose StatementBegin
CREATE TABLE fluxton.api_requests (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_uuid UUID NULL REFERENCES authentication.users(uuid),
    api_key UUID NULL,
    method VARCHAR(10) NOT NULL,
    endpoint TEXT NOT NULL,
    ip_address INET NOT NULL,
    user_agent TEXT NULL,
    params JSONB NULL,
    body JSONB NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE fluxton.api_requests;
-- +goose StatementEnd
