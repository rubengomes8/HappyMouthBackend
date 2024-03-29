package ingredients

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
	GetIngredients(ctx context.Context, reqOptions ReqOptions) ([]Ingredient, error)
}

type IngredientsHandler struct {
	svc      service
	tokenSvc corejwt.TokenService
}

func NewIngredientsHandler(svc service) IngredientsHandler {
	return IngredientsHandler{
		svc:      svc,
		tokenSvc: corejwt.NewTokenService(apiSecret, tokenLifespanHours),
	}
}

func (h IngredientsHandler) JWTAuthMiddleware() gin.HandlerFunc {
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

// GetIngredients is used to get the list of available ingredients
// ShowEntity godoc
// @tags ingredients
// @Summary Gets a list of ingredients.
// @Description Gets a list of ingredients.
// @Param sort-by-name query bool false "Sort ingredients by name."
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} []Ingredient
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/ingredients [get]
func (h IngredientsHandler) GetIngredients(ctx *gin.Context) {

	sortByName, err := getBoolQueryParamWithDefault(ctx, "sort-by-name", false)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	options := ReqOptions{
		SortByName: sortByName,
	}

	ingredients, err := h.svc.GetIngredients(ctx, options)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, ingredients)
	ctx.Writer.Flush()
}
