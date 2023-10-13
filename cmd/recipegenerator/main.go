package main

import (
	"log"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/rubengomes8/HappyMouthBackend/internal/recipegenerator"
)

const (
	kafkaBrokerAddress = "localhost:9092"
)

func main() {

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{kafkaBrokerAddress}, config)
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()

	recipeRouter := recipegenerator.NewAPI(producer)
	if err != nil {
		log.Fatal(err)
	}

	err = http.ListenAndServe(":8080", recipeRouter)
	if err != nil {
		log.Fatal(err)
	}
}
