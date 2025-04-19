package main

import (
	"log"
	"os"

	"github.com/AhmadAboElzahab/bridge/internal/database/seeder"
	"github.com/AhmadAboElzahab/bridge/internal/initializers"
)

func init() {
	initializers.LoadENV()
	initializers.ConnectDatabase()
}
func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run cmd/seed/main.go [countries|cities|all]")
	}

	seeder.Run(os.Args[1])
}
