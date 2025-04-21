-- +goose Up
-- +goose StatementBegin
CREATE TYPE field_type AS ENUM ('text', 'number', 'date', 'select', 'radio', 'checkbox');

CREATE TABLE fluxton.form_fields (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    form_uuid UUID NOT NULL,
    label VARCHAR(255) NOT NULL,
    type field_type NOT NULL,
    description TEXT NULL,
    is_required BOOLEAN DEFAULT FALSE,
    options JSONB NULL,
    min_length INT NULL,
    max_length INT NULL,
    min_value INT NULL,
    max_value INT NULL,
    pattern VARCHAR(255) NULL,
    default_value VARCHAR(255) NULL,
    start_date TIMESTAMP WITH TIME ZONE NULL,
    end_date TIMESTAMP WITH TIME ZONE NULL,
    date_format VARCHAR(50) NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    CONSTRAINT fk_form_fields_form FOREIGN KEY (form_uuid) REFERENCES fluxton.forms(uuid) ON DELETE CASCADE,
    CONSTRAINT unique_form_label UNIQUE (form_uuid, label)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS fluxton.form_fields CASCADE;
DROP TYPE IF EXISTS field_type;
-- +goose StatementEnd
