package recipes

import (
	"context"
	"time"

	"github.com/IBM/sarama"
	"github.com/go-resty/resty/v2"
)

const (
	gptRecipesTopic = "gpt-recipes"
	askGPT          = false
	useKafka        = false
	sleepTime       = 3
)

var (
	includeTemplate      = "Give me %s that includes the following ingredients: %s."
	excludeTemplate      = "Also, the recipe cannot have the following ingredients: %s."
	instructionsTemplate = "I would like to have only 4 sections separated by the pipe character |. Something like the following: name: x | ingredients: y | instructions: w | calories per serving: z. Also, split the list of ingredients by semicolon character ;"
)

//go:generate go-mockgen -f ./ -i iCache -d ./mocks/
type iCache interface {
	GetRecipeByKey(ctx context.Context, recipeKey string) (Recipe, error)
	GetRecipesByKeys(ctx context.Context, recipeKeys []string) ([]Recipe, error)
	StoreRecipe(ctx context.Context, recipeKey string, recipe Recipe) error
}

//go:generate go-mockgen -f ./ -i iUserRepo -d ./mocks/
type iUserRepo interface {
	GetUserRecipes(ctx context.Context, userID int) ([]UserRecipe, error)
	StoreUserRecipe(ctx context.Context, userRecipe UserRecipe) error
	UpdateUserRecipeFavorite(ctx context.Context, userID int, recipeKey string, isFavorite bool) error
}

// Implementing STRATEGY design pattern.
//
//go:generate go-mockgen -f ./ -i iRecipeAlgo -d ./mocks/
type iRecipeAlgo interface {
	BuildRecipe(ctx context.Context, definitions RecipeDefinitions) (string, error)
}

type Service struct {
	producer   sarama.SyncProducer
	cache      iCache
	userRepo   iUserRepo
	recipeAlgo iRecipeAlgo
}

func NewService(
	openAIEndpoint,
	openAIKey string,
	producer sarama.SyncProducer,
	cache iCache,
	userRepo iUserRepo,
) Service {

	var recipeAlgo iRecipeAlgo
	switch askGPT {
	case true:
		recipeAlgo = NewGPTRecipeAlgo(
			openAIEndpoint,
			openAIKey,
			resty.New(),
		)
	case false:
		recipeAlgo = fixedRecipeAlgo{}
	}

	return Service{
		producer:   producer,
		cache:      cache,
		userRepo:   userRepo,
		recipeAlgo: recipeAlgo,
	}
}

func (s Service) AskRecipe(ctx context.Context, recipeRequest RecipeDefinitions, userID int) (Recipe, error) {

	recipeKey := getRecipeKey(
		recipeRequest.RecipeType.Type,
		recipeRequest.IncludeIngredients,
		recipeRequest.ExcludeIngredients)

	recipe, err := s.cache.GetRecipeByKey(ctx, recipeKey)
	if err != nil {
		return Recipe{}, err
	}
	if recipe.HasTitle() {
		return recipe, nil
	}

	recipeStr, err := s.recipeAlgo.BuildRecipe(ctx, recipeRequest)
	if err != nil {
		return Recipe{}, err
	}

	if useKafka {
		_, _, err := ProduceRecipeEvent(s.producer, recipeKey, recipeStr, gptRecipesTopic)
		if err != nil {
			return Recipe{}, err
		}
	}

	parsedRecipe, err := parseRecipeString(recipeStr, recipeKey)
	if err != nil {
		return Recipe{}, err
	}

	parsedRecipe.Definitions = RecipeDefinitions{
		RecipeType:         recipeRequest.RecipeType,
		IncludeIngredients: recipeRequest.IncludeIngredients,
		ExcludeIngredients: recipeRequest.ExcludeIngredients,
	}

	err = s.cache.StoreRecipe(ctx, recipeKey, parsedRecipe)
	if err != nil {
		return Recipe{}, err
	}

	now := time.Now().UTC()
	err = s.userRepo.StoreUserRecipe(ctx, UserRecipe{
		UserID:    userID,
		RecipeKey: recipeKey,
		CreatedAt: &now,
		UpdatedAt: &now,
	})
	if err != nil {
		return Recipe{}, err
	}

	return parsedRecipe, nil
}

func (s Service) GetRecipesByUser(ctx context.Context, userID int) ([]Recipe, error) {

	userRecipes, err := s.userRepo.GetUserRecipes(ctx, userID)
	if err != nil {
		return []Recipe{}, err
	}

	var recipeKeys []string
	favoriteUserRecipes := map[string]struct{}{}
	for i := range userRecipes {
		if userRecipes[i].IsFavorite {
			favoriteUserRecipes[userRecipes[i].RecipeKey] = struct{}{}
		}
		recipeKeys = append(recipeKeys, userRecipes[i].RecipeKey)
	}

	recipes, err := s.cache.GetRecipesByKeys(ctx, recipeKeys)
	if err != nil {
		return []Recipe{}, err
	}

	for i := range recipes {
		if _, ok := favoriteUserRecipes[recipes[i].ID]; ok {
			recipes[i].IsFavorite = true
		}
	}

	return recipes, nil
}

func (s Service) SetUserRecipeFavorite(ctx context.Context, userID int, recipeKey string, isFavorite bool) error {
	return s.userRepo.UpdateUserRecipeFavorite(ctx, userID, recipeKey, isFavorite)
}
