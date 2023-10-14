package ingredients

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gorilla/mux"
)

type handler interface {
	GetIngredients(http.ResponseWriter, *http.Request)
}

type API struct {
	handler handler
}

func NewAPI() (*mux.Router, error) {

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == dynamodb.ServiceID && region == "us-east-1" {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           "http://localhost:8000",
				SigningRegion: "us-east-1",
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		return nil, err
	}

	dynamoDBClient := dynamodb.NewFromConfig(cfg)

	svc := NewService(dynamoDBClient)
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
