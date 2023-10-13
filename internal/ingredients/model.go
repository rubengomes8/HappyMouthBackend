package ingredients

import "github.com/gofrs/uuid"

type Ingredient struct {
	ID         uuid.UUID `json:"uuid" dynamodbav:"id"`
	Name       string    `json:"name" dynamodbav:"name"`
	IsFavorite bool      `json:"is_favorite" dynamodbav:"is_favorite"`
}
