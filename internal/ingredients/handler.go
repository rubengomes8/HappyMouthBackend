package ingredients

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:generate go-mockgen -f ./ -i service -d ./mocks/
type service interface {
	GetIngredients(ctx context.Context, reqOptions ReqOptions) ([]Ingredient, error)
}

type Handler struct {
	svc service
}

func NewHandler(svc service) Handler {
	return Handler{
		svc: svc,
	}
}

func (h Handler) GetIngredients(ctx *gin.Context) {

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
