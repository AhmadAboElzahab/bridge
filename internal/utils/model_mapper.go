package utils

import (
	"fmt"

	"github.com/AhmadAboElzahab/bridge/internal/models"
)

func GetModelForDataSource(modelName string) (interface{}, string, error) {
	switch modelName {
	case "Skill":
		return models.Skill{}, "skills", nil
	case "Language":
		return models.Language{}, "languages", nil
	case "Country":
		return models.Country{}, "countries", nil
	// Add more as needed
	default:
		return nil, "", fmt.Errorf("unsupported model: %s", modelName)
	}
}
