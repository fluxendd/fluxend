CREATE TABLE IF NOT EXISTS public.authors (
     id SERIAL PRIMARY KEY,
     name VARCHAR(255) NOT NULL,
     bio TEXT NOT NULL,
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO public.authors (name, bio) VALUES ('Jon Snow', 'The winter is coming.');
