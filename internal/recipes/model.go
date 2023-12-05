package recipes

import (
	"fmt"
	"strings"
	"time"

	"github.com/rubengomes8/HappyMouthBackend/internal/recipes/enums"
	"github.com/rubengomes8/HappyMouthBackend/pkg/utils"
)

type Recipe struct {
	ID           string            `json:"id"`
	Title        string            `json:"title"`
	IsFavorite   bool              `json:"is_favorite"`
	Definitions  RecipeDefinitions `json:"definitions"`
	Ingredients  []string          `json:"ingredients"`
	Instructions []string          `json:"instructions"`
	Calories     *float64          `json:"calories"`
	CreatedAt    *time.Time        `json:"created_at"`
	UpdatedAt    *time.Time        `json:"updated_at"`
	DeletedAt    *time.Time        `json:"deleted_at"`
}

func (r Recipe) HasTitle() bool {
	return r.Title != ""
}

type RecipeType struct {
	ID   int              `json:"id"`
	Type enums.RecipeType `json:"type"`
}

type RecipeDefinitions struct {
	IncludeIngredients []string   `json:"include_ingredients"`
	ExcludeIngredients []string   `json:"exclude_ingredients"`
	RecipeType         RecipeType `json:"recipe_type"`
}

func (r RecipeDefinitions) validate() error {

	if !r.RecipeType.Type.IsARecipeType() {
		return ErrInvalidRecipeType
	}

	if len(r.IncludeIngredients) == 0 {
		return ErrRequiredIncludeIngredients
	}

	mapUniqueIncludeIngredients := make(map[string]struct{})
	for _, includeIngredient := range r.IncludeIngredients {
		_, ok := mapUniqueIncludeIngredients[includeIngredient]
		if !ok {
			mapUniqueIncludeIngredients[includeIngredient] = struct{}{}
		}
	}

	mapUniqueExcludeIngredients := make(map[string]struct{})
	for _, excludeIngredient := range r.ExcludeIngredients {
		_, ok := mapUniqueExcludeIngredients[excludeIngredient]
		if !ok {
			mapUniqueExcludeIngredients[excludeIngredient] = struct{}{}
		}
	}

	for ingredient := range mapUniqueExcludeIngredients {
		_, ok := mapUniqueIncludeIngredients[ingredient]
		if ok {
			return ErrConflictingIngredients
		}
	}
	return nil
}

type RecipeFavoriteRequest struct {
	IsFavorite bool `json:"is_favorite"`
}

// key format: type|includedIngredients|excludedIngredients|timestamp.RFC3339
// key example: salad|mushroom,tomato|onion|2019-10-12T07:20:50.52Z
func getRecipeKey(
	recipeType enums.RecipeType,
	includeIngredients []string,
	excludeIngredients []string,
) string {
	recipeTypeLower := strings.ToLower(recipeType.String())
	uniqueIncludeSortedIngredients := utils.ToLowercaseUniqueSorted(includeIngredients)
	uniqueExcludeSortedIngredients := utils.ToLowercaseUniqueSorted(excludeIngredients)
	includeKey := strings.Join(uniqueIncludeSortedIngredients, ",")
	excludeKey := strings.Join(uniqueExcludeSortedIngredients, ",")
	return fmt.Sprintf("%s|%s|%s", recipeTypeLower, includeKey, excludeKey)
}

type UserRecipe struct {
	UserID     int        `json:"user_id,omitempty" gorm:"user_id"`
	RecipeKey  string     `json:"recipe_key,omitempty" gorm:"recipe_key"`
	IsFavorite bool       `json:"is_favorite,omitempty" gorm:"is_favorite"`
	CreatedAt  *time.Time `json:"created_at,omitempty" gorm:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty" gorm:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty" gorm:"deleted_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
