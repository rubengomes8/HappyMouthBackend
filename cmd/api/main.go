package main

import (
	"context"
	"log"

	"github.com/IBM/sarama"
	corepg "github.com/rubengomes8/HappyCore/pkg/postgres"
	"github.com/rubengomes8/HappyMouthBackend/internal/auth"
	"github.com/rubengomes8/HappyMouthBackend/internal/ingredients"
	"github.com/rubengomes8/HappyMouthBackend/internal/recipegenerator"
	"github.com/rubengomes8/HappyMouthBackend/pkg/kvstore"
	"github.com/rubengomes8/HappyMouthBackend/pkg/redis"
)

const (
	kafkaBrokerAddress = "localhost:9092"
	useKafka           = false
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

	// REDIS
	cache := redis.NewClient("localhost:6379", 0)
	if err := cache.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	// KAFKA
	var producer sarama.SyncProducer
	if useKafka {
		config := sarama.NewConfig()
		config.Producer.Return.Successes = true

		producer, err := sarama.NewSyncProducer([]string{kafkaBrokerAddress}, config)
		if err != nil {
			log.Fatal(err)
		}
		defer producer.Close()
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

	// RECIPE GENERATOR API - GIN + REDIS
	err = recipegenerator.NewAPI(cache, producer).Run("localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
}
