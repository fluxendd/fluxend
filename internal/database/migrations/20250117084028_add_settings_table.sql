-- +goose Up
-- +goose StatementBegin
CREATE TABLE fluxend.settings (
   id SERIAL PRIMARY KEY,
   name VARCHAR(255) NOT NULL UNIQUE,
   value VARCHAR(255) NOT NULL,
   default_value VARCHAR(255) NOT NULL,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE fluxend.settings;
-- +goose StatementEnd
