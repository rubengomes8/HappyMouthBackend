package recipes

import (
	"context"

	"github.com/rubengomes8/HappyMouthBackend/internal/users"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) userRepository {
	return userRepository{
		db: db,
	}
}

func (r userRepository) GetUserRecipes(ctx context.Context, userID int) ([]UserRecipe, error) {
	var userRecipes []UserRecipe
	err := r.db.WithContext(ctx).
		Model(UserRecipe{}).
		Where("user_id = ?", userID).
		Where("deleted_at IS NULL").
		Order("recipe_key, is_favorite").
		Scan(&userRecipes).
		Error
	if err != nil {
		return []UserRecipe{}, err
	}
	return userRecipes, nil

}

func (r userRepository) StoreUserRecipe(ctx context.Context, userRecipe UserRecipe) error {
	return r.db.WithContext(ctx).
		Create(userRecipe).
		Error
}

func (r userRepository) GetUserByUsername(ctx context.Context, username string) (users.User, error) {
	var user users.User
	err := r.db.WithContext(ctx).
		Model(users.User{}).
		Where("username = ?", username).
		Where("deleted_at IS NULL").
		First(&user).
		Error
	if err != nil {
		return users.User{}, err
	}
	return user, nil
}

func (r userRepository) UpdateUserRecipeFavorite(ctx context.Context, userID int, recipeKey string, isFavorite bool) error {
	return r.db.WithContext(ctx).
		Table("user_recipes").
		Where("user_id = ?", userID).
		Where("recipe_key = ?", recipeKey).
		Update("is_favorite", isFavorite).
		Error
}
