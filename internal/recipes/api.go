package recipes

import (
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/rubengomes8/HappyMouthBackend/pkg/redis"
	"gorm.io/gorm"
)

const (
	openAIEndpoint     = "https://api.openai.com/v1/chat/completions"
	openAIKey          = "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"
	kafkaBrokerAddress = "localhost:9092"
)

type Handler interface {
	JWTAuthMiddleware() gin.HandlerFunc
	CreateRecipe(ctx *gin.Context)
	GetRecipes(ctx *gin.Context)
}

type API struct {
	Handler Handler
}

func NewAPI(redis *redis.Cache, producer sarama.SyncProducer, db *gorm.DB) API {
	cache := NewCache(redis)
	userRepo := NewUserRepository(db)
	svc := NewService(openAIEndpoint, openAIKey, producer, cache, userRepo)
	h := NewRecipesHandler(svc)
	return API{
		Handler: h,
	}
}
