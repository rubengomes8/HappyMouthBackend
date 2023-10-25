package main

import (
	"encoding/csv"
	"fmt"
	"os"
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

	// 1. Abrir CSV
	// 2. Construir Mapa(Dicionário) map[string]string
	// 3. Fazer print do mapa
	// 4. Guardar filepaths dos JSON para processar
	// 5. Por cada ficheiro JSON, ler a receita e guardar numa BD

	filename := "recipes.csv" // Change this to the path of your CSV file

	data, err := readCSVFile(filename)
	if err != nil {
		fmt.Printf("Error reading CSV: %v\n", err)
		return
	}

	fmt.Println(data)
	for recipeKey, filePath := range data {
		fmt.Printf("%s: %s\n", recipeKey, filePath)
	}

}
