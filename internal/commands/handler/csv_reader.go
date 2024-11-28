package handler

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"vantsimulator/internal/models"
)

func Read(filePath string) ([]models.VANT, error) {
	vants := make([]models.VANT, 0)

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for _, vantElements := range records {

		if vantElements[0] == "id" {
			continue
		}

		id, _ := strconv.Atoi(vantElements[0])
		x, _ := strconv.ParseFloat(vantElements[1], 64)
		y, _ := strconv.ParseFloat(vantElements[2], 64)
		z, _ := strconv.ParseFloat(vantElements[3], 64)

		vant := models.VANT{
			ID: id,
			X:  x,
			Y:  y,
			Z:  z,
		}

		vants = append(vants, vant)
	}

	return vants, nil
}
