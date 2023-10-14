package ingredients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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

func (h Handler) GetIngredients(w http.ResponseWriter, r *http.Request) {

	sortByName, err := GetBoolQueryParamWithDefault(r, "sort-by-name", false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	options := ReqOptions{
		SortByName: sortByName,
	}

	ingredients, err := h.svc.GetIngredients(r.Context(), options)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to build recipe: %v", err).Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(ingredients)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to encode recipe response: %v", err).Error(), http.StatusInternalServerError)
	}
}
