package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	corejwt "github.com/rubengomes8/HappyCore/pkg/jwt"
	"github.com/rubengomes8/HappyMouthBackend/internal/users"
)

//go:generate go-mockgen -f ./ -i userService -d ./mocks/
type service interface {
	RegisterUser(ctx context.Context, user users.User) error
	LoginUser(ctx context.Context, req LoginInput) (string, error)
}

type Handler struct {
	svc service
}

func NewHandler(svc service) Handler {
	return Handler{
		svc: svc,
	}
}

func (h Handler) Register(ctx *gin.Context) {

	var input RegisterInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Sprintf("invalid body: %v", err))
		return
	}

	hashedPassword, err := corejwt.EncryptPassword(input.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	err = h.svc.RegisterUser(ctx, users.User{
		Username: input.Username,
		Passhash: hashedPassword,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("could not register user: %v", err))
		return
	}

	ctx.Writer.WriteHeader(http.StatusNoContent)
	ctx.Writer.Flush()
}

func (h Handler) Login(ctx *gin.Context) {

	var input LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Sprintf("invalid body: %v", err))
		return
	}

	token, err := h.svc.LoginUser(ctx, LoginInput{
		Username: input.Username,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("could not login user: %v", err))
		return
	}

	ctx.JSON(http.StatusOK, LoginResponse{Token: token})
	ctx.Writer.Flush()
}
