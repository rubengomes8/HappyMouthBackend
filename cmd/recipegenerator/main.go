package main

import (
	"context"
	"log"

	"github.com/IBM/sarama"
	"github.com/rubengomes8/HappyMouthBackend/internal/recipegenerator"
	"github.com/rubengomes8/HappyMouthBackend/pkg/redis"
)

const (
	kafkaBrokerAddress = "localhost:9092"
	useKafka           = false
)

func main() {

	ctx := context.Background()

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

	// REDIS
	cache := redis.NewClient("localhost:6379", 0)
	if err := cache.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	// API
	err := recipegenerator.NewAPI(cache, producer).Run("localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
}
