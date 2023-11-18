DROP INDEX idx_recipe_key_is_favorite;

ALTER TABLE user_recipes
DROP COLUMN is_favorite;