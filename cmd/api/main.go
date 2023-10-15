package main

import (
	"log"
	"net/http"

	"github.com/rubengomes8/HappyMouthBackend/internal/ingredients"
)

func main() {

	router, err := ingredients.NewAPI()
	if err != nil {
		log.Fatal(err)
	}

	err = http.ListenAndServe(":8082", router)
	if err != nil {
		log.Fatal(err)
	}
}
