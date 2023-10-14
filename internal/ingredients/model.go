package ingredients

type Ingredient struct {
	ID         string `json:"id" dynamodbav:"id"`
	Name       string `json:"name" dynamodbav:"name"`
	IsFavorite bool   `json:"is_favorite" dynamodbav:"is_favorite"`
}

type ReqOptions struct {
	SortByName bool
}
