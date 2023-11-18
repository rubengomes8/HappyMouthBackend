ALTER TABLE user_recipes
ADD COLUMN is_favorite BOOLEAN DEFAULT false;

CREATE INDEX idx_recipe_key_is_favorite ON user_recipes(recipe_key, is_favorite);
