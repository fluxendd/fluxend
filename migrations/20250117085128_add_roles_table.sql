-- +goose Up
-- +goose StatementBegin
CREATE TABLE roles (
   id SERIAL PRIMARY KEY,
   name VARCHAR(255) NOT NULL UNIQUE,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO roles (name) VALUES
/* Can do everything */
  ('king'),
/* can CRUD own organizations and projects underneath  */
  ('bishop'),
/* Editor: can CRUD projects underneath org he is part of but cannot modify org */
  ('lord'),
/* Explorer: can view projects underneath org he is part of */
  ('peasant');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE roles;
-- +goose StatementEnd
