package recipegenerator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getRecipeKey(t *testing.T) {
	type args struct {
		uniqueIncludeIngredients []string
		uniqueExcludeIngredients []string
	}
	tests := []struct {
		name      string
		args      args
		incSubKey string
		excSubKey string
	}{
		{
			name: "generate key",
			args: args{
				uniqueIncludeIngredients: []string{"tomato", "mushroom"},
				uniqueExcludeIngredients: []string{"onion", "honey"},
			},
			incSubKey: "mushroom,tomato",
			excSubKey: "honey,onion",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getRecipeKey(tt.args.uniqueIncludeIngredients, tt.args.uniqueExcludeIngredients)

			list := strings.Split(got, "|")

			assert.Equal(t, tt.incSubKey, list[0])
			assert.Equal(t, tt.excSubKey, list[1])
		})
	}
}
