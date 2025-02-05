-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS fluxton;
CREATE SCHEMA IF NOT EXISTS authentication;
CREATE SCHEMA IF NOT EXISTS storage;
CREATE SCHEMA IF NOT EXISTS public;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA IF EXISTS fluxton CASCADE;
DROP SCHEMA IF EXISTS authentication CASCADE;
DROP SCHEMA IF EXISTS storage CASCADE;
DROP SCHEMA IF EXISTS public CASCADE;
-- +goose StatementEnd
