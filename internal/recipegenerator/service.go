package recipegenerator

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/go-resty/resty/v2"
	"github.com/gofrs/uuid"

	"github.com/rubengomes8/HappyMouthBackend/internal/recipegenerator/examples"
	"github.com/rubengomes8/HappyMouthBackend/pkg/utils"
)

const (
	gptRecipesTopic = "gpt-recipes"
	askGPT          = false
	sleepTime       = 3
)

var (
	includeTemplate      = "Give me a simple recipe that includes the following ingredients: %s."
	excludeTemplate      = "Also, the recipe cannot have the following ingredients: %s."
	instructionsTemplate = "I would like to have only 4 sections separated by the pipe character |. Something like the following: name: x | ingredients: y | instructions: w | calories per serving: z. Also, split the list of ingredients by semicolon character ;"
)

type Service struct {
	OpenAIAPIEndpoint string
	OpenAIAPIKey      string
	OpenAIClient      *resty.Client
	Producer          sarama.SyncProducer
}

func NewService(openAIEndpoint, openAIKey string, producer sarama.SyncProducer) Service {
	return Service{
		OpenAIAPIEndpoint: openAIEndpoint,
		OpenAIAPIKey:      openAIKey,
		OpenAIClient:      resty.New(),
		Producer:          producer,
	}
}

func (s Service) AskRecipe(recipeRequest RecipeDefinitions) (Recipe, error) {

	chatGPTQuestion := createOpenAPIQuestion(
		utils.ToLowercaseUniqueSorted(recipeRequest.IncludeIngredients),
		utils.ToLowercaseUniqueSorted(recipeRequest.ExcludeIngredients))
	fmt.Println(chatGPTQuestion)

	var recipeStr string
	if askGPT {
		chatGPTResponse, err := s.getRecipeFromOpenAPI(chatGPTQuestion)
		if err != nil {
			return Recipe{}, err
		}

		body := chatGPTResponse.Body()

		var data map[string]interface{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			return Recipe{}, fmt.Errorf("error decoding JSON response: %v", err)
		}
		// fmt.Println(data)

		recipeStr, err = getOpenAPIRecipeString(data)
		if err != nil {
			return Recipe{}, err
		}

	} else {
		time.Sleep(sleepTime * time.Second)
		recipeStr = examples.Answer
	}

	eventKey := getRecipeKey(
		recipeRequest.IncludeIngredients,
		recipeRequest.ExcludeIngredients)
	_, _, err := ProduceRecipeEvent(s.Producer, eventKey, recipeStr, gptRecipesTopic)
	if err != nil {
		return Recipe{}, err
	}

	// fmt.Println("content")
	recipe, err := parseRecipeString(recipeStr)
	if err != nil {
		return Recipe{}, errors.New("error parsing recipe string")
	}

	// TODO: save recipe in cache or NoSQL database

	return recipe, nil
}

func createOpenAPIQuestion(includeIngredients, excludeIngredients []string) string {

	var question string
	var include string
	if len(includeIngredients) > 0 {
		include = strings.Join(includeIngredients, ", ")
		question += fmt.Sprintf(includeTemplate, include)
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

	response, err := s.OpenAIClient.R().
		SetAuthToken(s.OpenAIAPIKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"model": "gpt-3.5-turbo",
			"messages": []interface{}{map[string]interface{}{
				"role":    "system",
				"content": question}},
			"max_tokens": 500,
		}).
		Post(s.OpenAIAPIEndpoint)
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

func parseRecipeString(recipeStr string) (Recipe, error) {

	lowerRecipeStr := strings.ToLower(recipeStr)
	splittedByPipe := strings.Split(lowerRecipeStr, "|")

	if len(splittedByPipe) != 4 {
		return Recipe{}, errors.New("unexpected open api response")
	}

	now := time.Now().UTC()
	return Recipe{
		ID:           uuid.Must(uuid.NewV4()),
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
