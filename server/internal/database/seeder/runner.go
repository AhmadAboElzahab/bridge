package seeder

import (
	"log"
	"strings"

	"github.com/AhmadAboElzahab/bridge/internal/initializers"
	"github.com/AhmadAboElzahab/bridge/internal/models"
)

func Run(arg string) {
	initializers.ConnectDatabase()
	db := initializers.DB

	switch strings.ToLower(arg) {
	case "countries":
		SeedCountries()
	case "languages":
		SeedLanguages()
	case "skills":
		SeedSkills()
	case "maids":
		SeedMaidFormFields()
	case "users":
		SeedUserFormFields()
	case "userstabs":
		SeedUserTabs()
	// plop:cases
	case "fakemaids":
		var langs []models.Language
		var skills []models.Skill

		db.Find(&langs)
		db.Find(&skills)

		SeedMaids(db, langs, skills)
	case "all":
		SeedCountries()
		SeedLanguages()
		SeedSkills()
		SeedUserFormFields()
		SeedMaidFormFields()
		SeedUserTabs()
		// plop:all
	default:
		log.Fatalf("❌ Unknown seed type: %s\n", arg)
	}
}
