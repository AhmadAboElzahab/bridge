package seeder

import (
	"log"

	"github.com/AhmadAboElzahab/bridge/internal/initializers"
	"github.com/AhmadAboElzahab/bridge/internal/models"
)

func SeedUserFormFields() {
	fields := []models.FormField{
		{
			Label:        "First Name",
			FieldKey:     "FirstName",
			ModelName:    "User",
			FieldType:    "string",
			WidgetType:   "text",
			IsRequired:   true,
			IsFilterable: true,
			Order:        1,
		},
		{
			Label:        "Last Name",
			FieldKey:     "LastName",
			ModelName:    "User",
			FieldType:    "string",
			WidgetType:   "text",
			IsRequired:   true,
			IsFilterable: false,
			Order:        2,
		},
		{
			Label:        "Email",
			FieldKey:     "Email",
			ModelName:    "User",
			FieldType:    "string",
			WidgetType:   "email",
			IsRequired:   true,
			IsFilterable: true,
			Order:        3,
		},
		{
			Label:        "Date of Birth",
			FieldKey:     "DateOfBirth",
			ModelName:    "User",
			FieldType:    "date",
			WidgetType:   "date",
			IsRequired:   false,
			IsFilterable: false,
			Order:        4,
		},
		{
			Label:        "Avatar",
			FieldKey:     "Avatar",
			ModelName:    "User",
			FieldType:    "string",
			WidgetType:   "file",
			IsRequired:   false,
			IsFilterable: false,
			Order:        5,
		},
	}

	for _, field := range fields {
		if err := initializers.DB.Create(&field).Error; err != nil {
			log.Fatalf("failed to seed field: %v", err)
		}
	}
}
