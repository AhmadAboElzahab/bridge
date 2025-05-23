package utils

import (
	"fmt"

	"github.com/AhmadAboElzahab/bridge/internal/initializers"
	"github.com/AhmadAboElzahab/bridge/internal/models"
	"gorm.io/datatypes"
)

// CreateDefaultTabsForUserModel initializes a default "All" tab for a user and model,
// pre-populating it with column settings based on the model's FormField metadata.
//
// It inserts a new UserTab row with default filter/search settings,
// and adds corresponding UserTabColumn entries using each FormField's visibility,
// pinned state, and form width.
//
// Parameters:
//   - userID: ID of the user for whom the tab is being created
//   - model: name of the model (e.g., "Maid", "Driver") to create a tab for
func CreateDefaultTabsForUserModel(userID uint, model string) {
	var formFields []models.FormField
	initializers.DB.
		Where("model_name = ?", model).
		Order("field_order ASC").
		Find(&formFields)

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
		fmt.Printf("Creating column: FieldKey=%q, FormFieldID=%d\n", field.FieldKey, field.ID)

		column := models.UserTabColumn{
			UserTabID:   tab.ID,
			FormFieldID: &field.ID,
			FieldKey:    field.FieldKey,
			Visible:     field.TableIsVisible,
			Locked:      field.TableIsPinned,
			Order:       i + 1,
			Width:       field.FormWidth,
		}
		initializers.DB.Create(&column)
		fmt.Printf(">> Final FieldKey before insert: %q\n", column.FieldKey)
	}
}
