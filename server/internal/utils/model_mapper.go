package utils

import (
	"fmt"

	"github.com/AhmadAboElzahab/bridge/internal/models"
)

// GetModelForDataSource maps a model name (used in FormField.DataSource) to its corresponding
// Go struct instance and table name.
//
// This function enables dynamic resolution of models for use in query construction,
// form rendering, and search logic based on string model identifiers like "Skill", "Country", etc.
//
// Parameters:
//   - modelName: the name of the model as a string (e.g., "Skill", "Language", "Country")
//
// Returns:
//   - a zero-value instance of the model (e.g., models.Skill{})
//   - the corresponding table name as a string (e.g., "skills")
//   - an error if the model name is unsupported
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
