package main

import (
	"context"
	"log"
	"net/http"

	"github.com/rubengomes8/HappyMouthBackend/internal/ingredients"
	"github.com/rubengomes8/HappyMouthBackend/pkg/kvstore"
	"github.com/rubengomes8/HappyMouthBackend/pkg/redis"
)

func main() {

	ctx := context.Background()

	// REDIS
	cache := redis.NewClient("localhost:6379", 0)
	if err := cache.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	// DYNAMO DB
	dynamoDBClient, err := kvstore.NewClient(ctx, "http://localhost:8000", "us-east-1")
	if err := cache.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	router, err := ingredients.NewAPI(cache, dynamoDBClient)
	if err != nil {
		log.Fatal(err)
	}

	err = http.ListenAndServe(":8082", router)
	if err != nil {
		log.Fatal(err)
	}
}
