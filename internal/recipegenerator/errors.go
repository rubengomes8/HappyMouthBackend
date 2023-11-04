package recipegenerator

import "errors"

var (
	ErrRequiredIncludeIngredients = errors.New("happymouthbackend.recipes.error.including_ingredients_required")
	ErrConflictingIngredients     = errors.New("happymouthbackend.recipes.error.conflicting_ingredients")
	ErrInvalidUserID              = errors.New("happymouthbackend.recipes.error.invalid_user_id")
)
