CREATE TABLE user_recipes (
    user_recipe_id UUID PRIMARY KEY,
    recipe_key VARCHAR(255) NOT NULL,
    user_id INT NOT NULL REFERENCES users(id),
    created_at timestamp,
    updated_at timestamp,
    deleted_at timestamp
);

