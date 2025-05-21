CREATE TABLE IF NOT EXISTS public.posts (
     id SERIAL PRIMARY KEY,
     author_id INT NOT NULL REFERENCES public.authors(id),
     tag_id INT NOT NULL REFERENCES public.tags(id),
     title VARCHAR(255) NOT NULL,
     content TEXT NOT NULL,
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO posts (author_id, tag_id, title, content) VALUES (5, 12, 'Winter Is Coming', 'The cold winds rise, and with them come the shadows. The Starks must brace for the long night.');
INSERT INTO posts (author_id, tag_id, title, content) VALUES (14, 7, 'The Boy Who Lived', 'Even in the face of death, hope glimmers in the form of a lightning-shaped scar.');
INSERT INTO posts (author_id, tag_id, title, content) VALUES (2, 21, 'The Fellowship Forms', 'Nine walkers set out to destroy one ring that could doom them all.');
INSERT INTO posts (author_id, tag_id, title, content) VALUES (11, 4, 'Katniss Takes Aim', 'The girl on fire defies a corrupt Capitol with her arrow and her will.');
INSERT INTO posts (author_id, tag_id, title, content) VALUES (19, 33, 'Beneath the Misty Mountains', 'In the dark depths, a creature guards his precious, whispering riddles and madness.');
INSERT INTO posts (author_id, tag_id, title, content) VALUES (8, 1, 'The Maze Lies Ahead', 'Each wall a mystery, each turn a trap. Only the fastest survive the Glade.');
INSERT INTO posts (author_id, tag_id, title, content) VALUES (3, 28, 'The Da Vinci Code Unfolds', 'Symbols, secrets, and murder lead to revelations buried in art and faith.');
INSERT INTO posts (author_id, tag_id, title, content) VALUES (20, 3, 'Westeros Bleeds', 'Kings fall, bastards rise, and the game of thrones consumes all.');
INSERT INTO posts (author_id, tag_id, title, content) VALUES (7, 6, 'The Shire at Peace', 'Green hills, good ale, and second breakfasts are all a hobbit needs—until adventure knocks.');
INSERT INTO posts (author_id, tag_id, title, content) VALUES (1, 14, 'The Golden Compass Shines', 'In a world where daemons walk beside you, one girl seeks truth through dust.');
INSERT INTO posts (author_id, tag_id, title, content) VALUES (12, 25, 'Red Rising Begins', 'Born a slave among the mines, he rises to infiltrate the Golds and shatter their world.');
INSERT INTO posts (author_id, tag_id, title, content) VALUES (15, 17, 'The Hunger Games Ignite', 'When survival means killing, humanity is the price for rebellion.');
INSERT INTO posts (author_id, tag_id, title, content) VALUES (9, 11, 'Ender’s Game Starts', 'Trained to fight a war he barely understands, Ender becomes the weapon Earth needs.');
INSERT INTO posts (author_id, tag_id, title, content) VALUES (10, 9, 'Shadow and Bone Stir', 'In Ravka, darkness divides the nation—but a Sun Summoner may unite it.');
INSERT INTO posts (author_id, tag_id, title, content) VALUES (6, 31, 'American Gods Roam', 'Old deities walk among us, waging war with the new gods of technology and media.');
INSERT INTO posts (author_id, tag_id, title, content) VALUES (17, 10, 'Gone Girl Vanishes', 'A perfect marriage cracks when Amy disappears and the lies unravel fast.');
INSERT INTO posts (author_id, tag_id, title, content) VALUES (4, 20, 'The Name of the Wind', 'I have stolen princesses, burned down the town of Trebon, and spoken to gods.');
INSERT INTO posts (author_id, tag_id, title, content) VALUES (13, 13, 'Mistborn Rises', 'Ash falls from the sky, and rebellion brews under the Lord Ruler’s iron fist.');
INSERT INTO posts (author_id, tag_id, title, content) VALUES (16, 27, 'The Blade Itself Cuts Deep', 'Inquisitor Glokta tortures for the Union, but beneath the blood, truth festers.');
INSERT INTO posts (author_id, tag_id, title, content) VALUES (18, 2, 'Children of Time Evolve', 'As humanity falls, intelligent spiders rise to inherit the stars.');
INSERT INTO posts (author_id, tag_id, title, content) VALUES (3, 8, 'A Court of Thorns and Roses', 'Beneath the glamour, fae politics and ancient curses lie in wait for the unwary.');
INSERT INTO posts (author_id, tag_id, title, content) VALUES (6, 5, 'The Gunslinger Arrives', 'The man in black fled across the desert, and the gunslinger followed.');
INSERT INTO posts (author_id, tag_id, title, content) VALUES (17, 15, 'Dune: Fear is the Mind-Killer', 'He who controls the spice controls the universe—and Paul Atreides is its fulcrum.');
INSERT INTO posts (author_id, tag_id, title, content) VALUES (9, 32, 'It Waits in the Sewers', 'Derry’s children vanish, and something ancient awakens in the storm drains.');
INSERT INTO posts (author_id, tag_id, title, content) VALUES (13, 22, 'The Handmaid’s Lament', 'In Gilead, names are stripped, and women serve. Offred remembers her freedom.');
INSERT INTO posts (author_id, tag_id, title, content) VALUES (12, 1, 'The 100 Return', 'Earth is toxic, but survival demands they test it—and face what remains.');
INSERT INTO posts (author_id, tag_id, title, content) VALUES (7, 18, 'The Outsiders Stand Tall', 'Stay gold, Ponyboy, even when the world wants to crush you.');
INSERT INTO posts (author_id, tag_id, title, content) VALUES (2, 16, 'Circe Casts Her Spell', 'On an island of monsters and gods, a sorceress forges her destiny alone.');
INSERT INTO posts (author_id, tag_id, title, content) VALUES (15, 3, 'Fight Club’s First Rule', 'You do not talk about fight club. But rebellion always finds a voice.');
INSERT INTO posts (author_id, tag_id, title, content) VALUES (18, 29, 'The Martian Survives', 'I’m gonna have to science the shit out of this. Mars won’t kill me that easy.');
INSERT INTO posts (author_id, tag_id, title, content) VALUES (5, 19, 'The Night Circus Opens', 'A duel of magic, wrapped in illusion, played out in a circus only open at night.');
INSERT INTO posts (author_id, tag_id, title, content) VALUES (8, 6, 'Locke Lamora Lies Again', 'Thieves in Camorr don’t just steal money—they steal pride, and sometimes kingdoms.');
INSERT INTO posts (author_id, tag_id, title, content) VALUES (10, 27, 'The Ocean at the End of the Lane', 'Childhood memories and ancient power blur on a quiet English farm.');
INSERT INTO posts (author_id, tag_id, title, content) VALUES (20, 23, 'Project Hail Mary Begins', 'Earth is dying. Ryland Grace is its last hope—and he just woke up alone in space.');
INSERT INTO posts (author_id, tag_id, title, content) VALUES (11, 26, 'Thrawn Strikes Back', 'Cold logic, military brilliance, and blue skin make the Empire fear again.');
INSERT INTO posts (author_id, tag_id, title, content) VALUES (14, 2, 'The Pale King Watches', 'IRS agents dissect boredom and truth in a haunting unfinished epic.');
INSERT INTO posts (author_id, tag_id, title, content) VALUES (4, 11, 'Scythe Reaps Again', 'In a world without death, Scythes bring balance. Some enjoy it far too much.');
INSERT INTO posts (author_id, tag_id, title, content) VALUES (19, 10, 'Station Eleven Still Plays', 'The world ends, but Shakespeare remains, echoing in a post-apocalyptic silence.');
INSERT INTO posts (author_id, tag_id, title, content) VALUES (1, 17, 'Never Let Me Go', 'A love story tangled in clones, fate, and the quiet cruelty of acceptance.');
INSERT INTO posts (author_id, tag_id, title, content) VALUES (16, 30, 'The Power Awakens', 'Women seize the lightning, and the balance of the world flips—forever.');

