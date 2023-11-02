package auth

import (
	"context"

	"github.com/rubengomes8/HappyMouthBackend/internal/users"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	err := r.db.WithContext(ctx).
		Clauses(clause.Returning{}).
		Create(&user).
		Error
	if err != nil {
		return err
	}
	return nil
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
