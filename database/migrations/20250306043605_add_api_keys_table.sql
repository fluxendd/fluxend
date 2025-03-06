-- +goose Up
-- +goose StatementBegin
CREATE TABLE api_keys (
    id SERIAL PRIMARY KEY,
    key VARCHAR(255) NOT NULL,
    project_uuid UUID NOT NULL REFERENCES fluxton.projects(uuid),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE api_keys;
-- +goose StatementEnd
