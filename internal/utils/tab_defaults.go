package utils

import (
	"github.com/AhmadAboElzahab/bridge/internal/initializers"
	"github.com/AhmadAboElzahab/bridge/internal/models"
)

func CreateDefaultTabsForUserModel(userID uint, model string) {
	var defaults []models.DefaultUserTab
	initializers.DB.
		Where("model_name = ?", model).
		Find(&defaults)

	for _, d := range defaults {
		tab := models.UserTab{
			UserID:      userID,
			ModelName:   model,
			TabName:     d.TabName,
			IsDefault:   d.IsDefault,
			SearchTerm:  d.SearchTerm,
			Filters:     d.Filters,
			VisibleCols: d.VisibleCols,
		}
		initializers.DB.Create(&tab)
	}
}
