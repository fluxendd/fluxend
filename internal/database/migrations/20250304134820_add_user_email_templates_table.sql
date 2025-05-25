-- +goose Up
-- +goose StatementBegin
CREATE TABLE fluxend.user_email_templates (
    id SERIAL PRIMARY KEY,
    template_id INT NOT NULL REFERENCES fluxend.email_templates(id),
    user_uuid UUID NOT NULL REFERENCES authentication.users(uuid),
    message TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE fluxend.user_email_templates;
-- +goose StatementEnd
