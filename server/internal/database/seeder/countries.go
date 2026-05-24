package seeder

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/AhmadAboElzahab/bridge/internal/initializers"
	"github.com/AhmadAboElzahab/bridge/internal/models"
)

func SeedCountries() {
	EnsureDataFilesExist()

	file, err := os.Open("data/countryInfo.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) < 5 {
			continue
		}
		code := parts[0]
		name := parts[4]

		country := models.Country{Name: name, Code: code}
		initializers.DB.FirstOrCreate(&country, models.Country{Code: code})
		count++
	}
	log.Printf("✅ Seeded %d countries\n", count)
}
