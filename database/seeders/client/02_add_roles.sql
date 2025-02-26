CREATE ROLE web_anon NOINHERIT;

-- Grant SELECT (read) on all tables
GRANT SELECT ON ALL TABLES IN SCHEMA public TO web_anon;

-- Grant USAGE on the public schema
GRANT USAGE ON SCHEMA public TO web_anon;

-- Grant EXECUTE on all functions in the public schema (if needed)
GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA public TO web_anon;

-- Grant SELECT on all future tables in the public schema
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO web_anon;

-- Grant EXECUTE on all future functions in the public schema
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT EXECUTE ON FUNCTIONS TO web_anon;
