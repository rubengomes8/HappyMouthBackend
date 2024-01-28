package recipes

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
	AskRecipe(ctx context.Context, definitions RecipeDefinitions, userID int) (Recipe, error)
	GetRecipesByUser(ctx context.Context, userID int) ([]Recipe, error)
	SetUserRecipeFavorite(ctx context.Context, userID int, recipeKey string, isFavorite bool) error
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

// CreateRecipe is used to generate a new recipe.
// ShowEntity godoc
// @tags Recipes
// @Summary Generates a new recipe using OpenAI if it is a new set of parameters.
// @Description Generates a new recipe using OpenAI if it is a new set of parameters.
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body RecipeDefinitions true "Generate recipe request."
// @Success 200 {object} Recipe
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/recipes [post]
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

	userID, err := h.tokenSvc.ExtractClaimSub(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	recipe, err := h.svc.AskRecipe(ctx, recipeRequest, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, recipe)
	ctx.Writer.Flush()
}

// GetRecipes is used to get a list of recipes.
// ShowEntity godoc
// @tags Recipes
// @Summary Gets a list of recipes based on the provided filters.
// @Description Gets a list of recipes based on the provided filter
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} []Recipe
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/recipes [get]
func (h RecipesHandler) GetRecipes(ctx *gin.Context) {

	userID, err := h.tokenSvc.ExtractClaimSub(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	recipes, err := h.svc.GetRecipesByUser(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, recipes)
	ctx.Writer.Flush()
}

// SetUserRecipeFavorite is used to update the favorite state of a user recipe.
// ShowEntity godoc
// @tags Recipes
// @Summary Updates the favorite state of a user recipe.
// @Description Updates the favorite state of a user recipe.
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body RecipeDefinitions true "Update user recipe favorite request."
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/recipes/{id}/favorite [patch]
func (h RecipesHandler) SetUserRecipeFavorite(ctx *gin.Context) {

	userID, err := h.tokenSvc.ExtractClaimSub(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	recipeKey := ctx.Param("id")

	var req RecipeFavoriteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.svc.SetUserRecipeFavorite(ctx, userID, recipeKey, req.IsFavorite)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Writer.WriteHeader(http.StatusNoContent)
	ctx.Writer.Flush()
}
