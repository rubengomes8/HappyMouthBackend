package apiroutes

import (
	"github.com/gin-gonic/gin"
	"github.com/rubengomes8/HappyMouthBackend/internal/auth"
	"github.com/rubengomes8/HappyMouthBackend/internal/ingredients"
	"github.com/rubengomes8/HappyMouthBackend/internal/recipes"
)

func SetAPIRoutes(
	auth auth.Handler,
	recipes recipes.Handler,
	ingredients ingredients.Handler,
) *gin.Engine {

	r := gin.Default()

	// AUTH
	rAuthV1 := r.Group("/v1/auth")
	{
		rAuthV1.POST("/register", auth.Register)
		rAuthV1.POST("/login", auth.Login)
	}

	// INGREDIENTS
	v1Ingredients := r.Group("/v1/ingredients")
	{
		v1Ingredients.Use(ingredients.JWTAuthMiddleware())
		v1Ingredients.GET("", ingredients.GetIngredients)
	}

	// RECIPES
	v1Recipes := r.Group("/v1/recipes")
	{
		v1Recipes.Use(recipes.JWTAuthMiddleware())
		v1Recipes.POST("", recipes.CreateRecipe)
		v1Recipes.GET("", recipes.GetRecipes)
	}

	return r
}
