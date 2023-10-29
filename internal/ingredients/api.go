package ingredients

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	GetIngredients(ctx *gin.Context)
}

type API struct {
	Handler Handler
}

func NewAPI(dynamoDB *dynamodb.Client) API {
	repo := NewRepository(dynamoDB)
	svc := NewService(repo)
	h := NewIngredientsHandler(svc)
	return API{
		Handler: h,
	}
}
