CREATE TABLE IF NOT EXISTS public.authors (
     id SERIAL PRIMARY KEY,
     name VARCHAR(255) NOT NULL,
     bio TEXT NOT NULL,
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO public.authors (name, bio) VALUES
    ('Jon Snow', 'The winter is coming.'),
    ('Arya Stark', 'Faceless and fearless.'),
    ('Tyrion Lannister', 'I drink and I know things.'),
    ('Daenerys Targaryen', 'Breaker of chains, mother of dragons.'),
    ('Gandalf Greyhame', 'A wizard is never late.'),
    ('Frodo Baggins', 'One ring to rule them all.'),
    ('Sherlock Holmes', 'Consulting detective extraordinaire.'),
    ('Tony Stark', 'Genius, billionaire, playboy, philanthropist.'),
    ('Bruce Wayne', 'Vengeance wears a cape.'),
    ('Darth Vader', 'I am your father.'),
    ('Luke Skywalker', 'Farm boy turned Jedi.'),
    ('Leia Organa', 'A princess with a cause.'),
    ('Walter White', 'Say my name.'),
    ('Jesse Pinkman', 'Yeah science, bitch!'),
    ('Rick Grimes', 'This is not a democracy anymore.'),
    ('Homer Simpson', 'Mmm... donuts.'),
    ('Lisa Simpson', 'Saxophone prodigy and moral compass.'),
    ('Michael Scott', 'World’s best boss.'),
    ('Dwight Schrute', 'Bears, beets, Battlestar Galactica.'),
    ('Yoda', 'Do or do not, there is no try.');

