CREATE TABLE user_recipes (
    user_recipe_id SERIAL PRIMARY KEY,
    recipe_key VARCHAR(255) NOT NULL,
    user_id INT NOT NULL REFERENCES users(id)
);

