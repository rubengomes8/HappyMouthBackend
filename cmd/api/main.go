package main

import (
	"context"
	"log"
	"net/http"

	corepg "github.com/rubengomes8/HappyCore/pkg/postgres"
	"github.com/rubengomes8/HappyMouthBackend/internal/auth"
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

	// POSTGRES GORM DB
	conf := corepg.NewConfig("localhost", "database", "postgres", "passw0rd123", 5432)
	postgresDB, err := corepg.InitGormDB(conf)
	if err != nil {
		log.Fatal(err)
	}

	// AUTH API - GIN + GORM
	authAPI := auth.NewAPI(postgresDB)
	authRouter := authAPI.SetupRouter()
	err = authRouter.Run("localhost:8083")
	if err != nil {
		log.Fatal(err)
	}

	// INGREDIENTS ROUTER - GORILLA MUX
	ingredientsRouter, err := ingredients.NewAPI(dynamoDBClient)
	if err != nil {
		log.Fatal(err)
	}

	err = http.ListenAndServe(":8082", ingredientsRouter)
	if err != nil {
		log.Fatal(err)
	}
}
