package recipegenerator

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:generate go-mockgen -f ./ -i service -d ./mocks/
type service interface {
	AskRecipe(context.Context, RecipeDefinitions) (Recipe, error)
}

type Handler struct {
	svc service
}

func NewHandler(svc service) Handler {
	return Handler{
		svc: svc,
	}
}

func (h Handler) CreateRecipe(ctx *gin.Context) {

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
