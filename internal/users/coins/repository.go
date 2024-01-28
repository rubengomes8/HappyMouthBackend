package coins

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) repository {
	return repository{
		db: db,
	}
}

func (r repository) GetUserCoins(ctx context.Context, userID int) (UserCoins, error) {
	var userCoins UserCoins
	err := r.db.WithContext(ctx).
		Model(UserCoins{}).
		Where("user_id = ?", userID).
		Scan(&userCoins).
		Error
	if err != nil {
		return UserCoins{}, err
	}
	return userCoins, nil

}

func (r repository) UpsertUserCoin(ctx context.Context, userCoin UserCoins) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"coins", "updated_at"}),
		}).Create(&userCoin).
		Error
}
