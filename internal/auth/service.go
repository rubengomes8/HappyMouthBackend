package auth

import (
	"context"

	"github.com/rubengomes8/HappyMouthBackend/internal/users"

	corejwt "github.com/rubengomes8/HappyCore/pkg/jwt"
)

type repo interface {
	GetUserByUsername(ctx context.Context, username string) (users.User, error)
	CreateUser(ctx context.Context, user users.User) error
}

const (
	apiSecret = "86448213-7373-47B4-B3A2-55E4D8F1B987" // TODO: unsafe here
)

type Service struct {
	tokenSvc corejwt.TokenService
	repo     repo
}

func NewService(repo repo) Service {
	return Service{
		tokenSvc: corejwt.NewTokenService(apiSecret, 0),
		repo:     repo,
	}
}

func (s Service) LoginUser(ctx context.Context, req LoginInput) (string, error) {

	user, err := s.repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return "", err
	}

	err = corejwt.VerifyPassword(req.Password, user.Passhash)
	if err != nil {
		return "", err
	}

	token, err := s.tokenSvc.GenerateToken(uint(user.ID))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s Service) RegisterUser(ctx context.Context, user users.User) error {
	return s.repo.CreateUser(ctx, user)
}
