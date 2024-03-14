CREATE TABLE IF NOT EXISTS sold_recipes_history (
    id INTEGER PRIMARY KEY,
    recipe_id INTEGER NOT NULL,
    units INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY(recipe_id) REFERENCES recipe(id)
);