package recipe

import "errors"

var (
	ErrRequiredIncludeIngredients = errors.New("happymouthbackend.error.including_ingredients_required")
	ErrConflictingIngredients     = errors.New("happymouthbackend.error.conflicting_ingredients")
)
