CREATE TABLE IF NOT EXISTS public.products (
     id SERIAL PRIMARY KEY,
     title VARCHAR(255) NOT NULL,
     description TEXT NOT NULL,
     price DECIMAL(10, 2) NOT NULL,
     stock INT NOT NULL,
     has_discount BOOLEAN DEFAULT FALSE,
     total_sales INT DEFAULT 0,
     last_sold TIMESTAMP,
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO products (title, description, price, stock, has_discount, total_sales, last_sold) VALUES
  ('The Iron Blade', 'A legendary sword forged in dragon fire.', 129.99, 10, TRUE, 57, '2025-05-18 14:23:00'),
  ('Potion of Healing', 'Restores 50 HP instantly. Must-have for adventurers.', 9.99, 200, FALSE, 1243, '2025-05-20 09:11:00'),
  ('Enchanted Cloak', 'Grants temporary invisibility. Limited edition.', 249.50, 5, TRUE, 12, '2025-05-10 18:30:00'),
  ('Spell Tome: Fireball', 'Learn the devastating fireball spell.', 39.95, 50, FALSE, 326, '2025-05-19 16:45:00'),
  ('Leather Boots', 'Sturdy boots perfect for long treks.', 59.00, 75, FALSE, 145, '2025-05-21 12:00:00'),
  ('Elixir of Speed', 'Doubles movement speed for 5 minutes.', 19.99, 120, TRUE, 812, '2025-05-20 08:50:00'),
  ('Steel Shield', 'Heavy shield that blocks all frontal damage.', 89.99, 30, FALSE, 98, '2025-05-17 11:22:00'),
  ('Mana Crystal', 'Regenerates mana over time. Essential for mages.', 24.95, 60, FALSE, 563, '2025-05-21 07:35:00'),
  ('Mystic Ring', 'Increases intelligence and mana pool.', 149.99, 15, TRUE, 31, '2025-05-15 19:05:00'),
  ('Traveler’s Backpack', 'Carries twice the usual inventory.', 39.99, 90, FALSE, 274, '2025-05-19 13:55:00'),
  ('Dragon Scale Armor', 'Unmatched protection. Made from real dragon scales.', 399.99, 2, TRUE, 5, '2025-05-11 20:10:00'),
  ('Lantern of Guidance', 'Reveals hidden paths and traps.', 14.95, 40, FALSE, 412, '2025-05-21 09:20:00'),
  ('Rogue’s Dagger', 'High crit chance. Lightweight and deadly.', 74.99, 25, FALSE, 187, '2025-05-18 15:30:00'),
  ('Bag of Holding', 'A magical bag with near-infinite space.', 299.00, 8, TRUE, 14, '2025-05-14 17:40:00'),
  ('Phoenix Feather', 'Rare crafting material. Can revive once upon death.', 499.99, 1, TRUE, 2, '2025-05-09 22:00:00');
