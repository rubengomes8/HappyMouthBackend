package main

import (
	"context"
	"log"
	"net/http"

	"github.com/rubengomes8/HappyMouthBackend/internal/ingredients"
	"github.com/rubengomes8/HappyMouthBackend/pkg/kvstore"
)

func main() {

	ctx := context.Background()

	// DYNAMO DB
	dynamoDBClient, err := kvstore.NewClient(ctx, "http://localhost:8000", "us-east-1")
	if err != nil {
		log.Fatal(err)
	}

	router, err := ingredients.NewAPI(dynamoDBClient)
	if err != nil {
		log.Fatal(err)
	}

	err = http.ListenAndServe(":8082", router)
	if err != nil {
		log.Fatal(err)
	}
}
