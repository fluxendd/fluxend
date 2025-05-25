-- +goose Up
-- +goose StatementBegin
CREATE TABLE fluxend.api_logs (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_uuid UUID NULL,
    api_key UUID NULL,
    method VARCHAR(10) NOT NULL,
    status INTEGER NOT NULL,
    endpoint TEXT NOT NULL,
    ip_address INET NOT NULL,
    user_agent VARCHAR(255) NULL,
    params VARCHAR NULL,
    body TEXT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE fluxend.api_logs;
-- +goose StatementEnd
