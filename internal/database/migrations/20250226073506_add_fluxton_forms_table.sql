-- +goose Up
-- +goose StatementBegin
CREATE TABLE fluxend.forms (
   uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   project_uuid UUID NOT NULL,
   name VARCHAR(255) NOT NULL,
   description TEXT NULL,
   created_by UUID NOT NULL,
   updated_by UUID NULL,
   created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
   updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

   CONSTRAINT fk_forms_project FOREIGN KEY (project_uuid) REFERENCES fluxend.projects(uuid) ON DELETE CASCADE,
   CONSTRAINT fk_forms_created_by FOREIGN KEY (created_by) REFERENCES authentication.users(uuid) ON DELETE CASCADE,
   CONSTRAINT fk_forms_updated_by FOREIGN KEY (updated_by) REFERENCES authentication.users(uuid) ON DELETE SET NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS fluxend.forms CASCADE;
-- +goose StatementEnd
