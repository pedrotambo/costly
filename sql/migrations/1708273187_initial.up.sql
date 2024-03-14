CREATE TABLE IF NOT EXISTS ingredient (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    unit TEXT NOT NULL,
    price FLOAT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    last_modified TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS recipe (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL,
    last_modified TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS recipe_ingredient (
    recipe_id INTEGER NOT NULL,
    ingredient_id INTEGER NOT NULL,
    units INTEGER NOT NULL,
    PRIMARY KEY (recipe_id, ingredient_id),
    FOREIGN KEY(recipe_id) REFERENCES recipe(id),
    FOREIGN KEY(ingredient_id) REFERENCES ingredient(id)
);