-- +goose Up
-- +goose StatementBegin
DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'web_anon') THEN
            CREATE ROLE web_anon NOINHERIT;
        END IF;

        IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'authenticator') THEN
            CREATE ROLE authenticator NOINHERIT LOGIN PASSWORD 'fluxend';
            GRANT web_anon to authenticator;
        END IF;
    END $$;

CREATE SCHEMA IF NOT EXISTS fluxton;
CREATE SCHEMA IF NOT EXISTS authentication;
CREATE SCHEMA IF NOT EXISTS storage;

SET search_path TO public, fluxton, authentication, storage;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA IF EXISTS fluxton CASCADE;
DROP SCHEMA IF EXISTS authentication CASCADE;
DROP SCHEMA IF EXISTS storage CASCADE;
-- Avoid dropping the public schema as it's required for system tables
-- +goose StatementEnd
