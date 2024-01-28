package recipes

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/rubengomes8/HappyMouthBackend/internal/recipes/examples"
	"github.com/stretchr/testify/assert"
)

var (
	lowerRecipeStr = strings.ToLower(examples.Answer)
	splittedByPipe = strings.Split(lowerRecipeStr, "|")
)

func Test_parseRecipeName(t *testing.T) {

	type args struct {
		recipeName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "get recipe name",
			args: args{
				recipeName: splittedByPipe[0],
			},
			want: "tomato mushroom salad",
		},
	}
	for _, tt := range tests {
		fmt.Println(splittedByPipe[0])
		t.Run(tt.name, func(t *testing.T) {
			if got := parseRecipeName(tt.args.recipeName); got != tt.want {
				t.Errorf("getRecipeName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseRecipeCalories(t *testing.T) {
	type args struct {
		recipeCalories string
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "get recipe calories",
			args: args{
				recipeCalories: splittedByPipe[3],
			},
			want: 75.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseRecipeCalories(tt.args.recipeCalories); *got != tt.want {
				t.Errorf("getRecipeCalories() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseRecipeIngredients(t *testing.T) {
	type args struct {
		recipeIngredients string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "get recipe ingredients",
			args: args{
				recipeIngredients: splittedByPipe[1],
			},
			want: []string{
				"2 tomatoes",
				"1 cup mushrooms",
				"1 tablespoon olive oil",
				"1 tablespoon balsamic vinegar",
				"salt and pepper to taste",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseRecipeIngredients(tt.args.recipeIngredients); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getRecipeIngredients() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseRecipeInstructions(t *testing.T) {
	type args struct {
		recipeInstructions string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "get recipe instructions",
			args: args{
				recipeInstructions: splittedByPipe[2],
			},
			want: []string{
				"1. slice the tomatoes and mushrooms into bite-sized pieces.",
				"2. in a mixing bowl, combine the tomatoes and mushrooms.",
				"3. drizzle olive oil and balsamic vinegar over the mixture.",
				"4. season with salt and pepper to taste.",
				"5. gently toss the ingredients in the bowl until well combined.",
				"6. serve the tomato mushroom salad immediately.",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseRecipeInstructions(tt.args.recipeInstructions); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getRecipeInstructions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseRecipeString(t *testing.T) {
	type args struct {
		recipeStr string
		recipeKey string
	}
	tests := []struct {
		name    string
		args    args
		want    Recipe
		wantErr bool
	}{
		{
			name: "parse recipe string",
			args: args{
				recipeStr: examples.Answer,
				recipeKey: "tomato,mushroom",
			},
			want: Recipe{
				Title: "tomato mushroom salad",
				Ingredients: []string{
					"2 tomatoes",
					"1 cup mushrooms",
					"1 tablespoon olive oil",
					"1 tablespoon balsamic vinegar",
					"salt and pepper to taste",
				},
				Instructions: []string{
					"1. slice the tomatoes and mushrooms into bite-sized pieces.",
					"2. in a mixing bowl, combine the tomatoes and mushrooms.",
					"3. drizzle olive oil and balsamic vinegar over the mixture.",
					"4. season with salt and pepper to taste.",
					"5. gently toss the ingredients in the bowl until well combined.",
					"6. serve the tomato mushroom salad immediately.",
				},
				Calories: func() *float64 {
					f := 75.0
					return &f
				}(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseRecipeString(tt.args.recipeStr, tt.args.recipeKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRecipeString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want.Title, got.Title)
			assert.Equal(t, tt.want.Ingredients, got.Ingredients)
			assert.Equal(t, tt.want.Instructions, got.Instructions)
			assert.Equal(t, tt.want.Calories, got.Calories)
		})
	}
}
