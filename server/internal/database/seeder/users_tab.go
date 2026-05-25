package seeder

import (
	"log"

	"github.com/AhmadAboElzahab/bridge/internal/constants"
	"github.com/AhmadAboElzahab/bridge/internal/initializers"
	"github.com/AhmadAboElzahab/bridge/internal/models"
	"github.com/AhmadAboElzahab/bridge/internal/utils"
)

func SeedUserTabs() {
	db := initializers.DB

	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		log.Fatalf("❌ Failed to fetch users: %v", err)
	}

	if len(users) == 0 {
		log.Println("⚠️  No users found — skipping tab seed")
		return
	}

	for _, user := range users {
		for _, modelName := range []string{constants.USER, constants.MAID} {
			var existing models.UserTab
			err := db.Where("user_id = ? AND model_name = ? AND is_default = true", user.ID, modelName).First(&existing).Error
			if err == nil {
				continue
			}

			if err := utils.CreateDefaultTabsForUserModel(db, user.ID, modelName); err != nil {
				log.Fatalf("❌ Failed to create default %s tab for user %d: %v", modelName, user.ID, err)
			}
		}
	}

	log.Printf("✅ Seeded default tabs for %d users", len(users))
}
