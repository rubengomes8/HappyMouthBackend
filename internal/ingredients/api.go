package ingredients

import (
	"net/http"

	"github.com/gorilla/mux"
)

type handler interface {
	GetIngredients(http.ResponseWriter, *http.Request)
}

type API struct {
	handler handler
}

func NewAPI() *mux.Router {
	svc := NewService()
	h := NewHandler(svc)
	api := API{
		handler: h,
	}
	return api.SetIngredientsRoutes()
}

func (a API) SetIngredientsRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/ingredients", a.handler.GetIngredients).
		Methods(http.MethodGet)
	return r
}
