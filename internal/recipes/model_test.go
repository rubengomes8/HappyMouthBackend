package recipes

import (
	"strings"
	"testing"

	"github.com/rubengomes8/HappyMouthBackend/internal/recipes/enums"
	"github.com/stretchr/testify/assert"
)

func Test_getRecipeKey(t *testing.T) {
	type args struct {
		recipeType               enums.RecipeType
		uniqueIncludeIngredients []string
		uniqueExcludeIngredients []string
	}
	tests := []struct {
		name       string
		args       args
		incSubKey  string
		excSubKey  string
		typeSubKey string
	}{
		{
			name: "generate key",
			args: args{
				recipeType:               enums.Salad,
				uniqueIncludeIngredients: []string{"tomato", "mushroom"},
				uniqueExcludeIngredients: []string{"onion", "honey"},
			},
			typeSubKey: "salad",
			incSubKey:  "mushroom,tomato",
			excSubKey:  "honey,onion",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getRecipeKey(tt.args.recipeType, tt.args.uniqueIncludeIngredients, tt.args.uniqueExcludeIngredients)

			list := strings.Split(got, "|")

			assert.Equal(t, tt.typeSubKey, list[0])
			assert.Equal(t, tt.incSubKey, list[1])
			assert.Equal(t, tt.excSubKey, list[2])
		})
	}
}

func TestGenerateRecipeKey(t *testing.T) {
	recipeType := enums.Salad
	includedIngredients := []string{"mushroom", "tomato"}
	excludedIngredients := []string{"honey", "onion"}
	key := getRecipeKey(recipeType, includedIngredients, excludedIngredients)
	t.Logf("key: %v", key)
}
