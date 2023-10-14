package ingredients

type Ingredient struct {
	ID         string `json:"uuid" dynamodbav:"id"`
	Name       string `json:"name" dynamodbav:"name"`
	IsFavorite bool   `json:"is_favorite" dynamodbav:"is_favorite"`
}
