-- +goose Up
-- +goose StatementBegin
CREATE TABLE fluxton.form_field_responses (
  uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  form_response_uuid UUID NOT NULL,
  form_field_uuid UUID NOT NULL,
  value TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

  CONSTRAINT fk_field_responses_form_response FOREIGN KEY (form_response_uuid) REFERENCES fluxton.form_responses(uuid) ON DELETE CASCADE,
  CONSTRAINT fk_field_responses_form_field FOREIGN KEY (form_field_uuid) REFERENCES fluxton.form_fields(uuid) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS fluxton.form_field_responses CASCADE;
-- +goose StatementEnd
