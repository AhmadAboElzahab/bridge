package utils

import (
	"fmt"

	"github.com/AhmadAboElzahab/bridge/internal/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// CreateDefaultTabsForUserModel initializes a default "All" tab for a user and model,
// pre-populating it with column settings based on the model's FormField metadata.
func CreateDefaultTabsForUserModel(db *gorm.DB, userID uint, model string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var formFields []models.FormField
		if err := tx.Where("model_name = ?", model).Order("field_order ASC").Find(&formFields).Error; err != nil {
			return fmt.Errorf("failed to load form fields: %w", err)
		}

		tab := models.UserTab{
			UserID:     userID,
			ModelName:  model,
			TabName:    "All " + model + "s",
			IsDefault:  true,
			SearchTerm: "",
			Filters:    datatypes.JSON([]byte(`{}`)),
		}
		if err := tx.Create(&tab).Error; err != nil {
			return fmt.Errorf("failed to create default tab: %w", err)
		}

		columns := make([]models.UserTabColumn, 0, len(formFields))
		for i, field := range formFields {
			fieldID := field.ID
			columns = append(columns, models.UserTabColumn{
				UserTabID:   tab.ID,
				FormFieldID: &fieldID,
				FieldKey:    field.FieldKey,
				Visible:     field.TableIsVisible,
				Locked:      field.TableIsPinned,
				Order:       i + 1,
				Width:       field.FormWidth,
			})
		}

		if len(columns) > 0 {
			if err := tx.Create(&columns).Error; err != nil {
				return fmt.Errorf("failed to create default tab columns: %w", err)
			}
		}

		return nil
	})
}
