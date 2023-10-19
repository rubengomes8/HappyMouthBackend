package ingredients

import (
	"context"
	"sort"
)

type repo interface {
	GetIngredients(ctx context.Context) ([]Ingredient, error)
}

type Service struct {
	repo repo
}

func NewService(repo repo) Service {
	return Service{
		repo: repo,
	}
}

func (s Service) GetIngredients(ctx context.Context, reqOptions ReqOptions) ([]Ingredient, error) {

	ingredients, err := s.repo.GetIngredients(ctx)
	if err != nil {
		return []Ingredient{}, err
	}

	if reqOptions.SortByName {
		sort.Slice(ingredients, func(i, j int) bool {
			return ingredients[i].Name < ingredients[j].Name
		})
	}

	return ingredients, nil
}
