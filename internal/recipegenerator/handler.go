package recipegenerator

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//go:generate go-mockgen -f ./ -i service -d ./mocks/
type service interface {
	AskRecipe(RecipeDefinitions) (Recipe, error)
}

type Handler struct {
	svc service
}

func NewHandler(svc service) Handler {
	return Handler{
		svc: svc,
	}
}

func (h Handler) CreateRecipe(w http.ResponseWriter, r *http.Request) {

	var recipeRequest RecipeDefinitions
	err := json.NewDecoder(r.Body).Decode(&recipeRequest)
	if err != nil {
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	err = recipeRequest.validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	recipe, err := h.svc.AskRecipe(recipeRequest)
	if err != nil {
		http.Error(w, "failed to build recipe", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(recipe)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to encode recipe response: %v", err).Error(), http.StatusInternalServerError)
	}
}
