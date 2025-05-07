package seeder

import (
	"encoding/json"

	"github.com/AhmadAboElzahab/bridge/internal/constants"
	"github.com/AhmadAboElzahab/bridge/internal/initializers"
	"github.com/AhmadAboElzahab/bridge/internal/models"
	"gorm.io/datatypes"
)

func SeedDefaultUserTabs() error {
	defaults := []struct {
		ModelName   string
		TabName     string
		Search      string
		Filters     map[string]interface{}
		VisibleCols []string
	}{
		{
			ModelName:   constants.PATIENT,
			TabName:     "All Patients",
			Search:      "",
			Filters:     map[string]interface{}{},
			VisibleCols: []string{"id", "first_name", "country"},
		},
	}

	for _, d := range defaults {
		filtersJSON, _ := json.Marshal(d.Filters)
		colsJSON, _ := json.Marshal(d.VisibleCols)

		initializers.DB.FirstOrCreate(&models.DefaultUserTab{}, models.DefaultUserTab{
			ModelName:   d.ModelName,
			TabName:     d.TabName,
			IsDefault:   true,
			SearchTerm:  d.Search,
			Filters:     datatypes.JSON(filtersJSON),
			VisibleCols: datatypes.JSON(colsJSON),
		})
	}

	return nil
}
