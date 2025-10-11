CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT unique_user_category UNIQUE (user_id, name)
);

CREATE INDEX idx_categories_user_id ON categories(user_id);
