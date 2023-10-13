package ingredients

type Service struct {
}

func NewService() Service {
	return Service{}
}

func (s Service) GetIngredients() ([]Ingredient, error) {

	return []Ingredient{}, nil
}
