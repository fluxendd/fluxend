CREATE TABLE IF NOT EXISTS public.tags (
     id SERIAL PRIMARY KEY,
     name VARCHAR(255) NOT NULL,
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO public.tags (name) VALUES ('Travel');
INSERT INTO public.tags (name) VALUES ('Diary');
