package recipegenerator

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
	AskRecipe(context.Context, RecipeDefinitions) (Recipe, error)
}

type RecipesHandler struct {
	svc      service
	tokenSvc corejwt.TokenService
}

func NewRecipesHandler(svc service) RecipesHandler {
	return RecipesHandler{
		svc:      svc,
		tokenSvc: corejwt.NewTokenService(apiSecret, tokenLifespanHours),
	}
}

func (h RecipesHandler) JWTAuthMiddleware() gin.HandlerFunc {
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

func (h RecipesHandler) CreateRecipe(ctx *gin.Context) {

	var recipeRequest RecipeDefinitions
	if err := ctx.ShouldBindJSON(&recipeRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := recipeRequest.validate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recipe, err := h.svc.AskRecipe(ctx, recipeRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, recipe)
	ctx.Writer.Flush()
}

func (h RecipesHandler) GetRecipes(ctx *gin.Context) {
	panic("implement me")
}
