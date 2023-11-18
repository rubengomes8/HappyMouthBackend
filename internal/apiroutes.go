package apiroutes

import (
	"github.com/gin-gonic/gin"
	"github.com/rubengomes8/HappyMouthBackend/docs"
	"github.com/rubengomes8/HappyMouthBackend/internal/auth"
	"github.com/rubengomes8/HappyMouthBackend/internal/ingredients"
	"github.com/rubengomes8/HappyMouthBackend/internal/recipes"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

func SetAPIRoutes(
	auth auth.Handler,
	recipes recipes.Handler,
	ingredients ingredients.Handler,
) *gin.Engine {

	r := gin.Default()

	// SWAGGER
	docs.SwaggerInfo.Title = "Happy Mouth API"
	docs.SwaggerInfo.Description = "Happy Mouth REST API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080" // TODO: host should be set according to the environment
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
		v1Recipes.PATCH("/:id/favorite", recipes.SetUserRecipeFavorite)
	}

	return r
}
