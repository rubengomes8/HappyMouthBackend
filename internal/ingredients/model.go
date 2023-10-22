package ingredients

type Ingredient struct {
	ID   string `json:"id" dynamodbav:"id"`
	Name string `json:"name" dynamodbav:"name"`
}

type ReqOptions struct {
	SortByName bool
}
