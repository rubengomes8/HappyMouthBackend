package recipes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/IBM/sarama"
	random "github.com/Pallinder/go-randomdata"
	"github.com/go-resty/resty/v2"

	"github.com/rubengomes8/HappyMouthBackend/internal/recipes/enums"
	"github.com/rubengomes8/HappyMouthBackend/internal/recipes/examples"
	"github.com/rubengomes8/HappyMouthBackend/pkg/utils"
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

//go:generate go-mockgen -f ./ -i service -d ./mocks/
type iCache interface {
	GetRecipeByKey(ctx context.Context, recipeKey string) (Recipe, error)
	GetRecipesByKeys(ctx context.Context, recipeKeys []string) ([]Recipe, error)
	StoreRecipe(ctx context.Context, recipeKey string, recipe Recipe) error
}

//go:generate go-mockgen -f ./ -i service -d ./mocks/
type userRepo interface {
	GetUserRecipes(ctx context.Context, userID int) ([]UserRecipe, error)
	CreateUserRecipe(ctx context.Context, userRecipe UserRecipe) error
	UpdateUserRecipeFavorite(ctx context.Context, userID int, recipeKey string, isFavorite bool) error
}

type Service struct {
	openAIAPIEndpoint string
	openAIAPIKey      string
	openAIClient      *resty.Client
	producer          sarama.SyncProducer
	cache             iCache
	userRepo          userRepo
}

func NewService(
	openAIEndpoint,
	openAIKey string,
	producer sarama.SyncProducer,
	cache iCache,
	userRepo userRepo,
) Service {
	return Service{
		openAIAPIEndpoint: openAIEndpoint,
		openAIAPIKey:      openAIKey,
		openAIClient:      resty.New(),
		producer:          producer,
		cache:             cache,
		userRepo:          userRepo,
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

	var recipeStr string
	if askGPT {
		chatGPTQuestion := createOpenAPIQuestion(
			recipeRequest.RecipeType.Type,
			utils.ToLowercaseUniqueSorted(recipeRequest.IncludeIngredients),
			utils.ToLowercaseUniqueSorted(recipeRequest.ExcludeIngredients))
		fmt.Println(chatGPTQuestion)

		chatGPTResponse, err := s.getRecipeFromOpenAPI(chatGPTQuestion)
		if err != nil {
			return Recipe{}, err
		}

		body := chatGPTResponse.Body()

		var data map[string]interface{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			return Recipe{}, err
		}

		recipeStr, err = getOpenAPIRecipeString(data)
		if err != nil {
			return Recipe{}, err
		}

	} else {
		time.Sleep(sleepTime * time.Second)
		recipeStr = fmt.Sprintf(examples.Answer, random.SillyName())
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
	err = s.userRepo.CreateUserRecipe(ctx, UserRecipe{
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

func createOpenAPIQuestion(recipeType enums.RecipeType, includeIngredients, excludeIngredients []string) string {

	var question string
	var include string
	if len(includeIngredients) > 0 {
		include = strings.Join(includeIngredients, ", ")

		var recipeTypeString string
		if recipeType == enums.Any {
			recipeTypeString = "any recipe"
		} else {
			recipeTypeString = fmt.Sprintf("a simple %s recipe", strings.ToLower(recipeType.String()))
		}
		question += fmt.Sprintf(includeTemplate, recipeTypeString, include)
	}

	var exclude string
	if len(excludeIngredients) > 0 {
		exclude = strings.Join(excludeIngredients, ", ")
		question += fmt.Sprintf(" "+excludeTemplate, exclude)
	}

	question += " " + instructionsTemplate

	return question
}

func (s Service) getRecipeFromOpenAPI(question string) (*resty.Response, error) {

	response, err := s.openAIClient.R().
		SetAuthToken(s.openAIAPIKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"model": "gpt-3.5-turbo",
			"messages": []interface{}{map[string]interface{}{
				"role":    "system",
				"content": question}},
			"max_tokens": 500,
		}).
		Post(s.openAIAPIEndpoint)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func getOpenAPIRecipeString(data map[string]interface{}) (string, error) {

	_, ok := data["error"]
	if ok {
		errorCode := data["error"].(map[string]interface{})["code"].(string)
		return "", fmt.Errorf("open ai error: %v", errorCode)
	}

	_, ok = data["choices"]
	if !ok {
		return "", errors.New("response have no field choices")
	}

	choices := data["choices"].([]interface{})
	if len(choices) == 0 {
		return "", errors.New("response choices are empty")
	}

	choice := choices[0].(map[string]interface{})
	_, ok = choice["message"]
	if !ok {
		return "", errors.New("choice have no field message")
	}

	message := choice["message"].(map[string]interface{})
	_, ok = message["content"]
	if !ok {
		return "", errors.New("choice message have no field content")
	}

	return message["content"].(string), nil
}

func parseRecipeString(recipeStr, recipeKey string) (Recipe, error) {

	lowerRecipeStr := strings.ToLower(recipeStr)
	splittedByPipe := strings.Split(lowerRecipeStr, "|")

	if len(splittedByPipe) != 4 {
		return Recipe{}, errors.New("unexpected open api response")
	}

	now := time.Now().UTC()
	return Recipe{
		ID:           recipeKey,
		Title:        getRecipeName(splittedByPipe[0]),
		Ingredients:  getRecipeIngredients(splittedByPipe[1]),
		Instructions: getRecipeInstructions(splittedByPipe[2]),
		Calories:     getRecipeCalories(splittedByPipe[3]),
		CreatedAt:    &now,
		UpdatedAt:    &now,
	}, nil
}

func getRecipeName(recipeName string) string {
	splittedByName := strings.Split(recipeName, "name:")
	return strings.TrimSpace(splittedByName[len(splittedByName)-1])
}

func getRecipeCalories(recipeCalories string) *float64 {
	splittedByColon := strings.Split(recipeCalories, ":")
	caloriesStr := strings.TrimSpace(splittedByColon[len(splittedByColon)-1])
	calories, err := strconv.ParseFloat(caloriesStr, 32)
	if err != nil {
		return nil
	}
	return &calories
}

func getRecipeIngredients(recipeIngredients string) []string {
	splittedByColon := strings.Split(recipeIngredients, ":")
	splittedBySemicolon := strings.Split(splittedByColon[len(splittedByColon)-1], ";")
	var ingredients []string
	for i := range splittedBySemicolon {
		ingredients = append(ingredients, strings.TrimSpace(splittedBySemicolon[i]))
	}
	return ingredients
}

func getRecipeInstructions(recipeInstructions string) []string {
	splittedByColon := strings.Split(recipeInstructions, ": ")
	splittedByNewline := strings.Split(splittedByColon[len(splittedByColon)-1], "\n")
	var instructions []string
	for i := range splittedByNewline {
		if splittedByNewline[i] == "\n" || splittedByNewline[i] == " " || splittedByNewline[i] == "" {
			continue
		}
		instructions = append(instructions, strings.TrimSpace(splittedByNewline[i]))
	}
	return instructions
}
