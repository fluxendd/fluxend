CREATE TABLE IF NOT EXISTS posts (
     id SERIAL PRIMARY KEY,
     title VARCHAR(255) NOT NULL,
     content TEXT NOT NULL,
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO posts (title, content) VALUES ('Hello World', 'This is a test post.');
INSERT INTO posts (title, content) VALUES ('Hello World 2', 'This is another test post.');
