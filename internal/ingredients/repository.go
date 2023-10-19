package ingredients

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type repository struct {
	db *dynamodb.Client
}

func NewRepository(db *dynamodb.Client) repository {
	return repository{
		db: db,
	}
}

func (r repository) GetIngredients(ctx context.Context) ([]Ingredient, error) {

	input := &dynamodb.ScanInput{
		TableName: aws.String("ingredients"),
	}

	result, err := r.db.Scan(ctx, input)
	if err != nil {
		return []Ingredient{}, err
	}

	var ingredients []Ingredient
	for _, item := range result.Items {
		var ingredient Ingredient
		err := attributevalue.UnmarshalMap(item, &ingredient)
		if err != nil {
			return []Ingredient{}, err
		}
		ingredients = append(ingredients, ingredient)
	}

	return ingredients, nil
}
