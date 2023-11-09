package recipes

import (
	"context"

	"github.com/rubengomes8/HappyMouthBackend/pkg/redis"
)

type cache struct {
	redis *redis.Cache
}

func NewCache(redis *redis.Cache) cache {
	return cache{
		redis: redis,
	}
}

func (c cache) GetRecipeByKey(ctx context.Context, key string) (Recipe, error) {
	var recipe Recipe
	err := c.redis.Get(ctx, key, &recipe)
	if err != nil && err != redis.ErrNotFound {
		return Recipe{}, err
	}
	return recipe, nil
}

func (c cache) GetRecipesByKeys(ctx context.Context, recipeKeys []string) ([]Recipe, error) {
	var recipes []Recipe
	// TODO: should be nice to do a single redis call, instead of N.
	for _, recipeKey := range recipeKeys {
		recipe, err := c.GetRecipeByKey(ctx, recipeKey)
		if err != nil {
			return []Recipe{}, err
		}
		recipes = append(recipes, recipe)
	}
	return recipes, nil
}

func (c cache) StoreRecipe(ctx context.Context, key string, recipe Recipe) error {
	return c.redis.Set(ctx, key, recipe, 0 /* ttl */)
}
