package ingredients

import (
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gorilla/mux"
)

type handler interface {
	GetIngredients(http.ResponseWriter, *http.Request)
}

type API struct {
	handler handler
}

func NewAPI(dynamoDB *dynamodb.Client) (*mux.Router, error) {

	repo := NewRepository(dynamoDB)
	svc := NewService(repo)
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
