-- +goose Up
-- +goose StatementBegin
CREATE TABLE roles (
   id SERIAL PRIMARY KEY,
   name VARCHAR(255) NOT NULL UNIQUE,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO roles (name) VALUES
  ('king'), /* Super admin: and can do anything. Can create/edit/delete any org or project */
  ('bishop'), /* Admin: can create/edit/delete own orgs and projects  */
  ('lord'), /* Editor: Can create/edit/delete projects under an org */
  ('peasant'); /* Explorer: readonly access to everything under an org */

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE roles;
-- +goose StatementEnd
