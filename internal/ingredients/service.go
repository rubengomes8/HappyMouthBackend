package ingredients

import (
	"context"

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

func (s Service) GetIngredients(ctx context.Context) ([]Ingredient, error) {

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

	return ingredients, nil
}
