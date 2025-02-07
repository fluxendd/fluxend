-- +goose Up
-- +goose StatementBegin
CREATE TABLE authentication.users (
   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   role_id INT NOT NULL REFERENCES authentication.roles (id),
   username VARCHAR(255) NOT NULL UNIQUE,
   email VARCHAR(255) NOT NULL UNIQUE,
   status VARCHAR(10) NOT NULL CHECK (status IN ('active', 'inactive')),
   password VARCHAR(255) NOT NULL,
   bio TEXT DEFAULT '',
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_users_email ON authentication.users(email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE authentication.users;
-- +goose StatementEnd
