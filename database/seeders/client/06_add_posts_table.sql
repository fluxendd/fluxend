CREATE TABLE IF NOT EXISTS public.posts (
     id SERIAL PRIMARY KEY,
     author_id INT NOT NULL REFERENCES public.authors(id),
     tag_id INT NOT NULL REFERENCES public.tags(id),
     title VARCHAR(255) NOT NULL,
     content TEXT NOT NULL,
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO posts (author_id, tag_id, title, content) VALUES (1, 1, 'My first post', 'This is the content of my first post.');
INSERT INTO posts (author_id, tag_id, title, content) VALUES (1, 1, 'My second post', 'This is the content of my second post.');
