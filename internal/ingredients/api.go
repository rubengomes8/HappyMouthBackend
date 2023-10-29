package ingredients

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
)

type handler interface {
	GetIngredients(ctx *gin.Context)
}

type API struct {
	handler handler
}

func NewAPI(dynamoDB *dynamodb.Client) *gin.Engine {

	repo := NewRepository(dynamoDB)
	svc := NewService(repo)
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
		v1.GET("/ingredients", a.handler.GetIngredients)
	}
	return r
}
