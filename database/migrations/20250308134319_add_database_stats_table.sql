-- +goose Up
-- +goose StatementBegin
CREATE TABLE database_stats (
    id SERIAL PRIMARY KEY,
    database_name VARCHAR(255) NOT NULL,
    total_size VARCHAR(255) NOT NULL,
    index_size VARCHAR(255) NOT NULL,
    unused_index JSONB NOT NULL,
    table_count JSONB NOT NULL,
    table_size JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE database_stats;
-- +goose StatementEnd
