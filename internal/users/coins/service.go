package coins

import (
	"context"
)

//go:generate go-mockgen -f ./ -i iRepo -d ./mocks/
type iRepo interface {
	GetUserCoins(ctx context.Context, userID int) (UserCoins, error)
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
