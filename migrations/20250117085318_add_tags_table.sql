-- +goose Up
-- +goose StatementBegin
CREATE TABLE tags (
   id SERIAL PRIMARY KEY,
   name VARCHAR(255) NOT NULL,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tags;
-- +goose StatementEnd
