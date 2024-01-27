package recipes

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	random "github.com/Pallinder/go-randomdata"

	"github.com/go-resty/resty/v2"
	"github.com/rubengomes8/HappyMouthBackend/internal/recipes/examples"
	"github.com/rubengomes8/HappyMouthBackend/pkg/utils"
)

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

type fixedRecipeAlgo struct {
}

func (a fixedRecipeAlgo) BuildRecipe(ctx context.Context, definitions RecipeDefinitions) (string, error) {
	time.Sleep(sleepTime * time.Second)
	return fmt.Sprintf(examples.Answer, random.SillyName()), nil
}
