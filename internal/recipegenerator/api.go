package recipegenerator

import (
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/rubengomes8/HappyMouthBackend/pkg/redis"
)

const (
	openAIEndpoint     = "https://api.openai.com/v1/chat/completions"
	openAIKey          = "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"
	kafkaBrokerAddress = "localhost:9092"
)

type Handler interface {
	JWTAuthMiddleware() gin.HandlerFunc
	CreateRecipe(ctx *gin.Context)
}

type API struct {
	Handler Handler
}

func NewAPI(cache *redis.Cache, producer sarama.SyncProducer) API {
	repo := NewRepository(cache)
	svc := NewService(openAIEndpoint, openAIKey, producer, repo)
	h := NewRecipesHandler(svc)
	return API{
		Handler: h,
	}
}
