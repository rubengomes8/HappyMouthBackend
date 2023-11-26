package auth

import (
	"context"
	"time"

	"github.com/rubengomes8/HappyMouthBackend/internal/users"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{
		db: db,
	}
}

func (r Repository) CreateUser(ctx context.Context, user users.User) error {
	return r.db.WithContext(ctx).
		Create(&user).
		Error
}

func (r Repository) GetUserByUsername(ctx context.Context, username string) (users.User, error) {
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

func (r Repository) GetUserByID(ctx context.Context, userID int) (users.User, error) {
	var user users.User
	err := r.db.WithContext(ctx).
		Model(users.User{}).
		Where("id = ?", userID).
		Where("deleted_at IS NULL").
		First(&user).
		Error
	if err != nil {
		return users.User{}, err
	}
	return user, nil
}

func (r Repository) UpdatePassword(ctx context.Context, username string, passhash string) error {
	now := time.Now().UTC()
	return r.db.WithContext(ctx).
		Model(users.User{}).
		Where("username = ?", username).
		Updates(map[string]interface{}{
			"passhash":   passhash,
			"updated_at": now,
		}).
		Error
}
