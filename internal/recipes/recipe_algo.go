package recipes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	random "github.com/Pallinder/go-randomdata"

	"github.com/go-resty/resty/v2"
	"github.com/rubengomes8/HappyMouthBackend/internal/recipes/enums"
	"github.com/rubengomes8/HappyMouthBackend/internal/recipes/examples"
	"github.com/rubengomes8/HappyMouthBackend/pkg/utils"
)

// GPT Recipe algo

type gptRecipeAlgo struct {
	endpoint string
	apiKey   string
	client   *resty.Client
}

func NewGPTRecipeAlgo(
	endpoint string,
	apiKey string,
	client *resty.Client,
) gptRecipeAlgo {
	return gptRecipeAlgo{
		endpoint: endpoint,
		apiKey:   apiKey,
		client:   client,
	}
}

func (a gptRecipeAlgo) getRecipeFromOpenAPI(question string) (*resty.Response, error) {

	response, err := a.client.R().
		SetAuthToken(a.apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"model": "gpt-3.5-turbo",
			"messages": []interface{}{map[string]interface{}{
				"role":    "system",
				"content": question}},
			"max_tokens": 500,
		}).
		Post(a.endpoint)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (a gptRecipeAlgo) BuildRecipe(ctx context.Context, definitions RecipeDefinitions) (string, error) {

	var recipeStr string
	chatGPTQuestion := createOpenAPIQuestion(
		definitions.RecipeType.Type,
		utils.ToLowercaseUniqueSorted(definitions.IncludeIngredients),
		utils.ToLowercaseUniqueSorted(definitions.ExcludeIngredients))
	fmt.Println(chatGPTQuestion)

	chatGPTResponse, err := a.getRecipeFromOpenAPI(chatGPTQuestion)
	if err != nil {
		return "", err
	}

	body := chatGPTResponse.Body()

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}

	recipeStr, err = getOpenAPIRecipeString(data)
	if err != nil {
		return "", err
	}

	return recipeStr, nil
}

func createOpenAPIQuestion(recipeType enums.RecipeType, includeIngredients, excludeIngredients []string) string {

	var question string
	var include string
	if len(includeIngredients) > 0 {
		include = strings.Join(includeIngredients, ", ")

		var recipeTypeString string
		if recipeType == enums.Any {
			recipeTypeString = "a simple recipe"
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
		Title:        parseRecipeName(splittedByPipe[0]),
		Ingredients:  parseRecipeIngredients(splittedByPipe[1]),
		Instructions: parseRecipeInstructions(splittedByPipe[2]),
		Calories:     parseRecipeCalories(splittedByPipe[3]),
		CreatedAt:    &now,
		UpdatedAt:    &now,
	}, nil
}

func parseRecipeName(recipeName string) string {
	splittedByName := strings.Split(recipeName, "name:")
	return strings.TrimSpace(splittedByName[len(splittedByName)-1])
}

func parseRecipeCalories(recipeCalories string) *float64 {
	splittedByColon := strings.Split(recipeCalories, ":")
	caloriesStr := strings.TrimSpace(splittedByColon[len(splittedByColon)-1])
	calories, err := strconv.ParseFloat(caloriesStr, 32)
	if err != nil {
		return nil
	}
	return &calories
}

func parseRecipeIngredients(recipeIngredients string) []string {
	splittedByColon := strings.Split(recipeIngredients, ":")
	splittedBySemicolon := strings.Split(splittedByColon[len(splittedByColon)-1], ";")
	var ingredients []string
	for i := range splittedBySemicolon {
		ingredients = append(ingredients, strings.TrimSpace(splittedBySemicolon[i]))
	}
	return ingredients
}

func parseRecipeInstructions(recipeInstructions string) []string {
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

// Fake Recipe algo

type fixedRecipeAlgo struct {
}

func (a fixedRecipeAlgo) BuildRecipe(ctx context.Context, definitions RecipeDefinitions) (string, error) {
	time.Sleep(sleepTime * time.Second)
	return fmt.Sprintf(examples.Answer, random.SillyName()), nil
}
