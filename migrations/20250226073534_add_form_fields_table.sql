-- +goose Up
-- +goose StatementBegin
CREATE TABLE fluxton.form_fields (
     uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
     form_uuid UUID NOT NULL,
     label VARCHAR(255) NOT NULL,
     type TEXT CHECK (type IN ('text', 'number', 'date', 'select', 'radio', 'checkbox')) NOT NULL,
     is_required BOOLEAN DEFAULT FALSE,
     options JSONB NULL,
     created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
     updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

     CONSTRAINT fk_form_fields_form FOREIGN KEY (form_uuid) REFERENCES fluxton.forms(uuid) ON DELETE CASCADE,
     CONSTRAINT unique_form_label UNIQUE (form_uuid, label)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS fluxton.form_fields CASCADE;
-- +goose StatementEnd
