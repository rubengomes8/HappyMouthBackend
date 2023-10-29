package apiroutes

import (
	"github.com/gin-gonic/gin"
	"github.com/rubengomes8/HappyMouthBackend/internal/auth"
	"github.com/rubengomes8/HappyMouthBackend/internal/ingredients"
	"github.com/rubengomes8/HappyMouthBackend/internal/recipegenerator"
)

func SetAPIRoutes(
	auth auth.Handler,
	recipes recipegenerator.Handler,
	ingredients ingredients.Handler,
) *gin.Engine {

	r := gin.Default()
	v1 := r.Group("/v1")
	{
		// AUTH
		v1.POST("/auth/register", auth.Register)
		v1.POST("/auth/login", auth.Login)

		// INGREDIENTS
		v1.GET("/ingredients", ingredients.GetIngredients)

		// RECIPES
		v1.POST("/recipes", recipes.CreateRecipe)
	}

	return r
}
