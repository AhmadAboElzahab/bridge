package utils

import (
	"fmt"

	"github.com/AhmadAboElzahab/bridge/internal/initializers"
	"github.com/AhmadAboElzahab/bridge/internal/models"
	"gorm.io/datatypes"
)

func CreateDefaultTabsForUserModel(userID uint, model string) {
	var formFields []models.FormField
	initializers.DB.Where("model_name = ?", model).Order("field_order ASC").Find(&formFields)

	tab := models.UserTab{
		UserID:     userID,
		ModelName:  model,
		TabName:    "All " + model + "s",
		IsDefault:  true,
		SearchTerm: "",
		Filters:    datatypes.JSON([]byte(`{}`)),
	}
	initializers.DB.Create(&tab)

	for i, field := range formFields {
		column := models.UserTabColumn{
			UserTabID:   tab.ID,
			FormFieldID: field.ID,
			Visible:     field.TableIsVisible,
			Locked:      field.TableIsPinned,
			Order:       i + 1,
			Width:       field.FormWidth,
		}
		initializers.DB.Create(&column)
		fmt.Printf("Created UserTabColumn for field: %s\n", field.FieldKey)
	}
}
