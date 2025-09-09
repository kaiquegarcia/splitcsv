package main

import "github.com/kaiquegarcia/splitcsv/internal/models"

func main() {
	input := models.NewInput()
	writer := models.NewWriter(input)
	writer.Split()
}
