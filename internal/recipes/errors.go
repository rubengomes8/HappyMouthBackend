package recipes

import "errors"

var (
	ErrRequiredIncludeIngredients = errors.New("happymouthbackend.recipes.error.including_ingredients_required")
	ErrConflictingIngredients     = errors.New("happymouthbackend.recipes.error.conflicting_ingredients")
	ErrRequiredParam              = errors.New("happymouthbackend.recipes.error.required_param")
	ErrInvalidInt                 = errors.New("happymouthbackend.recipes.error.invalid_int")
	ErrInvalidRecipeType          = errors.New("happymouthbackend.recipes.error.invalid_type")
)
