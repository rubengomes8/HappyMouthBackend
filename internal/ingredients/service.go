package ingredients

import (
	"context"
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Service struct {
	DynamoDBClient *dynamodb.Client
}

func NewService(dynamoDBClient *dynamodb.Client) Service {
	return Service{
		DynamoDBClient: dynamoDBClient,
	}
}

func (s Service) GetIngredients(ctx context.Context, reqOptions ReqOptions) ([]Ingredient, error) {

	input := &dynamodb.ScanInput{
		TableName: aws.String("ingredients"),
	}

	result, err := s.DynamoDBClient.Scan(ctx, input)
	if err != nil {
		return []Ingredient{}, err
	}

	var ingredients []Ingredient
	for _, dynamoItem := range result.Items {

		var ingredient Ingredient

		err := attributevalue.UnmarshalMap(dynamoItem, &ingredient)
		if err != nil {
			return []Ingredient{}, err
		}

		ingredients = append(ingredients, ingredient)
	}

	if reqOptions.SortByName {
		sort.Slice(ingredients, func(i, j int) bool {
			return ingredients[i].Name < ingredients[j].Name
		})
	}

	return ingredients, nil
}
