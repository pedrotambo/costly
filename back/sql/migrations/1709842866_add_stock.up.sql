ALTER TABLE ingredient
ADD units_in_stock INTEGER NOT NULL
DEFAULT 0;

CREATE TABLE IF NOT EXISTS stock_history (
    id INTEGER PRIMARY KEY,
    ingredient_id INTEGER NOT NULL,
    units INTEGER NOT NULL,
    price FLOAT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY(ingredient_id) REFERENCES ingredient(id)
);