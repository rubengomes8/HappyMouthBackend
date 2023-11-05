package recipes

import (
	"context"

	"github.com/rubengomes8/HappyMouthBackend/pkg/redis"
)

type repository struct {
	cache *redis.Cache
}

func NewRepository(cache *redis.Cache) repository {
	return repository{
		cache: cache,
	}
}

func (r repository) GetRecipeByKey(ctx context.Context, key string) (Recipe, error) {
	var recipe Recipe
	err := r.cache.Get(ctx, key, &recipe)
	if err != nil && err != redis.ErrNotFound {
		return Recipe{}, err
	}
	return recipe, nil
}

func (r repository) GetRecipesByKeys(ctx context.Context, recipeKey string) ([]Recipe, error) {
	panic("implement me")
}

func (r repository) StoreRecipe(ctx context.Context, key string, recipe Recipe) error {
	return r.cache.Set(ctx, key, recipe, 0 /* ttl */)
}
