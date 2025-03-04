-- +goose Up
-- +goose StatementBegin
CREATE TABLE fluxton.user_email_templates (
    id SERIAL PRIMARY KEY,
    template_id INT NOT NULL REFERENCES fluxton.email_templates(id),
    user_uuid UUID NOT NULL REFERENCES authentication.users(uuid),
    message TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE fluxton.user_email_templates;
-- +goose StatementEnd
