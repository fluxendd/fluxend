-- +goose Up
-- +goose StatementBegin
CREATE TABLE notes (
    id SERIAL PRIMARY KEY,
    /*tag_id INT NOT NULL,*/
    user_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE notes ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE notes;
-- +goose StatementEnd
