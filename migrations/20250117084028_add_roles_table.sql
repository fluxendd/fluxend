-- +goose Up
-- +goose StatementBegin
CREATE TABLE fluxton.settings (
   id SERIAL PRIMARY KEY,
   name VARCHAR(255) NOT NULL UNIQUE,
   value VARCHAR(255) NOT NULL,
   default_value VARCHAR(255) NOT NULL,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO fluxton.settings (name, value, default_value)
VALUES
    ('title', 'Fluxton', 'Fluxton'),
    ('description', 'Fluxton is a BAAS', 'Fluxton is a BAAS'),
    ('url', 'https://fluxton.com', 'https://fluxton.com'),
    ('allow_registration', 'true', 'true'),
    ('max_projects_per_user', '10', '10');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE fluxton.settings;
-- +goose StatementEnd
