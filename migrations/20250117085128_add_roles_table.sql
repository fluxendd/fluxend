-- +goose Up
-- +goose StatementBegin
CREATE TABLE authentication.roles (
   id SERIAL PRIMARY KEY,
   name VARCHAR(255) NOT NULL UNIQUE,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO authentication.roles (name) VALUES
/* Fluxton superuser with admin area access */
  ('superman'),
/* Can do everything with created orgs, project and users */
  ('owner'),
/* can CRUD own organizations and projects underneath  */
  ('admin'),
/* Editor: can CRUD projects underneath org he is part of but cannot modify org */
  ('developer'),
/* Explorer: can view projects underneath org he is part of */
  ('explorer');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE authentication.roles;
-- +goose StatementEnd
