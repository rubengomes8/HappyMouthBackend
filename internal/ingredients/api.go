package ingredients

import (
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gorilla/mux"
	"github.com/rubengomes8/HappyMouthBackend/pkg/redis"
)

type handler interface {
	GetIngredients(http.ResponseWriter, *http.Request)
}

type API struct {
	handler handler
}

func NewAPI(cache *redis.Cache, dynDB *dynamodb.Client) (*mux.Router, error) {

	svc := NewService(cache, dynDB)
	h := NewHandler(svc)
	api := API{
		handler: h,
	}

	return api.SetIngredientsRoutes(), nil
}

func (a API) SetIngredientsRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/ingredients", a.handler.GetIngredients).
		Methods(http.MethodGet)
	return r
}
