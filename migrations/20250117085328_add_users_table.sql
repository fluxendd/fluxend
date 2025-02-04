-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   role_id INT NOT NULL,
   username VARCHAR(255) NOT NULL UNIQUE,
   email VARCHAR(255) NOT NULL UNIQUE,
   status VARCHAR(10) NOT NULL CHECK (status IN ('active', 'inactive')),
   password VARCHAR(255) NOT NULL,
   bio TEXT DEFAULT '',
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_users_email ON users(email);
ALTER TABLE users ADD CONSTRAINT fk_users_role_id FOREIGN KEY (role_id) REFERENCES roles(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
