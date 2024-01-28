package coins

import (
	"context"

	"gorm.io/gorm"
)

//go:generate go-mockgen -f ./ -i iRepo -d ./mocks/
type iRepo interface {
	RunTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error
	GetUserCoins(ctx context.Context, userID int) (UserCoins, error)
	GetUserCoinsTx(tx *gorm.DB, userID int) (UserCoins, error)
	UpsertUserCoinTx(tx *gorm.DB, userCoins UserCoins) error
}

type Service struct {
	repo iRepo
}

func NewService(
	repo iRepo,
) Service {
	return Service{
		repo: repo,
	}
}

func (s Service) GetUserCoins(ctx context.Context, userID int) (UserCoins, error) {
	return s.repo.GetUserCoins(ctx, userID)
}

func (s Service) SubtractUserCoins(ctx context.Context, userID int, quantity int) error {

	// SQL transaction
	return s.repo.RunTransaction(ctx, func(tx *gorm.DB) error {

		userCoins, err := s.repo.GetUserCoinsTx(tx, userID)
		if err != nil {
			return err
		}

		newUserCoins, err := userCoins.Subtract(quantity)
		if err != nil {
			return err
		}

		err = s.repo.UpsertUserCoinTx(tx, newUserCoins)
		if err != nil {
			return err
		}

		return nil
	})
}
