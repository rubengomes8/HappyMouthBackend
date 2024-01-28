package users

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	corejwt "github.com/rubengomes8/HappyCore/pkg/jwt"
)

const (
	apiSecret          = "86448213-7373-47B4-B3A2-55E4D8F1B987" // TODO: unsafe here
	tokenLifespanHours = 8760
)

//go:generate go-mockgen -f ./ -i service -d ./mocks/
type service interface {
	GetUserCoins(ctx context.Context, userID int) (UserCoins, error)
}

type Handler struct {
	svc      service
	tokenSvc corejwt.TokenService
}

func NewHandler(svc service) Handler {
	return Handler{
		svc:      svc,
		tokenSvc: corejwt.NewTokenService(apiSecret, tokenLifespanHours),
	}
}

func (h Handler) JWTAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := h.tokenSvc.ValidateToken(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

type UserCoinsResponse struct {
	UserID int64 `json:"user_id"`
	Coins  int   `json:"coins"`
}

// GetCoins is used to get the user numebr of coins.
// ShowEntity godoc
// @tags users
// @Summary Gets the number of user coins.
// @Description Gets the number of user coins.
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} []UserCoinsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/users/coins [get]
func (h Handler) GetUserCoins(ctx *gin.Context) {

	userID, err := h.tokenSvc.ExtractClaimSub(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	userCoins, err := h.svc.GetUserCoins(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, UserCoinsResponse{
		UserID: userCoins.UserID,
		Coins:  userCoins.Coins,
	})
	ctx.Writer.Flush()
}
