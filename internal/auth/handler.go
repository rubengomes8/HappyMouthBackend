package auth

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	corejwt "github.com/rubengomes8/HappyCore/pkg/jwt"
	"github.com/rubengomes8/HappyMouthBackend/internal/users"
	"golang.org/x/crypto/bcrypt"
)

//go:generate go-mockgen -f ./ -i userService -d ./mocks/
type service interface {
	RegisterUser(ctx context.Context, user users.User) error
	LoginUser(ctx context.Context, req LoginInput) (string, error)
}

type AuthHandler struct {
	svc service
}

func NewAuthHandler(svc service) AuthHandler {
	return AuthHandler{
		svc: svc,
	}
}

func (h AuthHandler) Register(ctx *gin.Context) {

	var input RegisterInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := corejwt.EncryptPassword(input.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.svc.RegisterUser(ctx, users.User{
		Username: input.Username,
		Passhash: hashedPassword,
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
