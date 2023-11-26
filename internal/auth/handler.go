package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rubengomes8/HappyMouthBackend/internal/users"
	"golang.org/x/crypto/bcrypt"

	corejwt "github.com/rubengomes8/HappyCore/pkg/jwt"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

const (
	passwordMinEntropyBits = 50
)

//go:generate go-mockgen -f ./ -i userService -d ./mocks/
type service interface {
	RegisterUser(ctx context.Context, user users.User) error
	LoginUser(ctx context.Context, req LoginInput) (string, error)
	ChangePassword(ctx context.Context, req ChangePasswordRequest) error
}

type AuthHandler struct {
	svc      service
	tokenSvc corejwt.TokenService
}

func NewAuthHandler(svc service) AuthHandler {
	return AuthHandler{
		svc:      svc,
		tokenSvc: corejwt.NewTokenService(apiSecret, 99999 /* tokenLifespanHours */),
	}
}

func (h AuthHandler) Register(ctx *gin.Context) {

	var input RegisterInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := passwordvalidator.Validate(input.Password, passwordMinEntropyBits)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": ErrWeakPassword.Error()})
		return
	}

	hashedPassword, err := corejwt.EncryptPassword(input.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	now := time.Now().UTC()
	err = h.svc.RegisterUser(ctx, users.User{
		ID:        0,
		Username:  input.Username,
		Passhash:  hashedPassword,
		Email:     input.Email,
		CreatedAt: &now,
		UpdatedAt: &now,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Writer.WriteHeader(http.StatusNoContent)
	ctx.Writer.Flush()
}

func (h AuthHandler) Login(ctx *gin.Context) {

	var input LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.svc.LoginUser(ctx, LoginInput{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, LoginResponse{Token: token})
	ctx.Writer.Flush()
}

func (h AuthHandler) ChangePassword(ctx *gin.Context) {

	var req ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.svc.ChangePassword(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Writer.WriteHeader(http.StatusNoContent)
	ctx.Writer.Flush()
}
