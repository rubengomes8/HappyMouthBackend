package main

import (
	"context"
	"log"

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

	// AUTH API - GIN + PGSQL GORM
	err = auth.NewAPI(postgresDB).Run("localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	// INGREDIENTS ROUTER - GIN + DYNAMO
	err = ingredients.NewAPI(dynamoDBClient).Run("localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
}
