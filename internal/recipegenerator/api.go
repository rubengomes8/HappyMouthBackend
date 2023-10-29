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

type handler interface {
	CreateRecipe(ctx *gin.Context)
}

type API struct {
	handler handler
}

func NewAPI(cache *redis.Cache, producer sarama.SyncProducer) *gin.Engine {
	repo := NewRepository(cache)
	svc := NewService(openAIEndpoint, openAIKey, producer, repo)
	h := NewHandler(svc)
	api := API{
		handler: h,
	}
	return api.SetupRouter()
}

func (a API) SetupRouter() *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/v1")
	{
		v1.POST("/recipes", a.handler.CreateRecipe)
	}
	return r
}
