package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadENV() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
}
