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
	panic("implement me")
}

func (r userRepository) GetUserByUsername(ctx context.Context, username string) (users.User, error) {
	var user users.User
	err := r.db.WithContext(ctx).
		Model(users.User{}).
		Where("username = ?", username).
		First(&user).
		Error
	if err != nil {
		return users.User{}, err
	}
	return user, nil
}
