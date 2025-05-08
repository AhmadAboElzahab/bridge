package seeder

import (
	"log"
	"strings"

	"github.com/AhmadAboElzahab/bridge/internal/initializers"
)

func Run(arg string) {
	// Connect to the database
	initializers.ConnectDatabase()

	switch strings.ToLower(arg) {
	case "countries":
		SeedCountries()
	case "languages":
		SeedLanguages()
	case "skills":
		SeedSkills()
	case "all":
		SeedCountries()
		SeedUserFormFields()
		SeedPatientFormFields()
	default:
		log.Fatalf("❌ Unknown seed type: %s\n", arg)
	}
}
