package recipegenerator

import (
	"net/http"

	"github.com/IBM/sarama"
	"github.com/gorilla/mux"
)

const (
	openAIEndpoint     = "https://api.openai.com/v1/chat/completions"
	openAIKey          = "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"
	kafkaBrokerAddress = "localhost:9092"
)

type handler interface {
	CreateRecipe(http.ResponseWriter, *http.Request)
}

type API struct {
	handler handler
}

func NewAPI(producer sarama.SyncProducer) *mux.Router {

	svc := NewService(openAIEndpoint, openAIKey, producer)
	h := NewHandler(svc)
	api := API{
		handler: h,
	}
	return api.SetRoutes()
}

func (a API) SetRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/recipes", a.handler.CreateRecipe).
		Methods(http.MethodPost)
	return r
}