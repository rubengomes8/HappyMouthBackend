package main

import (
	"context"
	"log"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/rubengomes8/HappyMouthBackend/internal/recipegenerator"
	"github.com/rubengomes8/HappyMouthBackend/pkg/redis"
)

const (
	kafkaBrokerAddress = "localhost:9092"
)

func main() {

	ctx := context.Background()

	// KAFKA
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{kafkaBrokerAddress}, config)
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()

	// REDIS
	cache := redis.NewClient("localhost:6379", 0)
	if err := cache.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	// API
	recipeRouter := recipegenerator.NewAPI(cache, producer)
	if err != nil {
		log.Fatal(err)
	}

	err = http.ListenAndServe(":8080", recipeRouter)
	if err != nil {
		log.Fatal(err)
	}
}
