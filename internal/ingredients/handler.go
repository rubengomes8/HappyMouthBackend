package ingredients

import "net/http"

//go:generate go-mockgen -f ./ -i service -d ./mocks/
type service interface {
}

type Handler struct {
	svc service
}

func NewHandler(svc service) Handler {
	return Handler{
		svc: svc,
	}
}

func (h Handler) GetIngredients(w http.ResponseWriter, r *http.Request) {}
