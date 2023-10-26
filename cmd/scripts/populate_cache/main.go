package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/rubengomes8/HappyMouthBackend/internal/recipegenerator"
	"github.com/rubengomes8/HappyMouthBackend/pkg/redis"
)

const (
	pathToPopulateCacheDir = "cmd/scripts/populate_cache/"
)

func readCSVFile(filename string) (map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'

	data := make(map[string]string)

	// Read and parse the CSV data
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		if len(record) == 2 {
			recipeKey := record[0]
			filePath := record[1]
			data[recipeKey] = filePath
		}
	}

	return data, nil
}

func main() {

	ctx := context.Background()

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	filename := "recipes.csv"
	if !strings.Contains(wd, "cmd") {
		filename = pathToPopulateCacheDir + filename
	}

	mapRecipeKeyToFilePath, err := readCSVFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	cache := redis.NewClient("localhost:6379", 0)
	if err := cache.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	var fileRecipe recipegenerator.Recipe
	for recipeKey, filePath := range mapRecipeKeyToFilePath {
		if recipeKey == "recipe_key" {
			continue
		}

		if !strings.Contains(wd, "cmd") {
			filePath = pathToPopulateCacheDir + filePath
		}

		file, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal([]byte(file), &fileRecipe)
		if err != nil {
			log.Fatal(err)
		}

		err = cache.Set(ctx, recipeKey, fileRecipe, 0)
		if err != nil {
			log.Fatal(err)
		}
	}
}
