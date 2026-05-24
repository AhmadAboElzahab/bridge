package seeder

import (
	"log"

	"github.com/AhmadAboElzahab/bridge/internal/constants"
	"github.com/AhmadAboElzahab/bridge/internal/initializers"
	"github.com/AhmadAboElzahab/bridge/internal/models"
)

func SeedUserFormFields() {
	fields := []models.FormField{
		{
			Label:     "First Name",
			FieldKey:  "FirstName",
			ModelName: constants.USER,
			HelpText:  "",

			FormFieldType:  constants.TypeStringField,
			FormIsRequired: true,
			FormOrder:      1,
			FormStage:      "info",
			FormWidth:      2,

			TableIsPinned:  false,
			TableIsVisible: true,
			TableOrder:     4,
		},
		{
			Label:     "Last Name",
			FieldKey:  "LastName",
			ModelName: constants.USER,
			HelpText:  "",

			FormFieldType:  constants.TypeStringField,
			FormWidth:      2,
			FormOrder:      1,
			FormStage:      "info",
			FormIsRequired: true,

			TableIsPinned:  false,
			TableIsVisible: true,
			TableOrder:     4,
		},
		{
			Label:     "Email",
			FieldKey:  "Email",
			ModelName: constants.USER,
			HelpText:  "",

			FormFieldType:  constants.TypeEmailField,
			FormWidth:      2,
			FormOrder:      1,
			FormStage:      "info",
			FormIsRequired: true,

			TableIsPinned:  false,
			TableIsVisible: true,
			TableOrder:     4,
		},
		{
			Label:     "Date of Birth",
			FieldKey:  "DateOfBirth",
			ModelName: constants.USER,
			HelpText:  "",

			FormFieldType:  constants.TypeDateField,
			FormIsRequired: false,
			FormOrder:      4,
			FormStage:      "info",
			FormWidth:      2,

			TableIsPinned:  false,
			TableIsVisible: true,
			TableOrder:     4,
		},
		{
			Label:     "Avatar",
			FieldKey:  "Avatar",
			ModelName: constants.USER,
			HelpText:  "",

			FormFieldType:  constants.TypeImageField,
			FormIsRequired: false,
			FormOrder:      5,
			FormStage:      "info",
			FormWidth:      2,

			TableIsPinned:  false,
			TableIsVisible: true,
			TableOrder:     0,
		},
	}

	for _, field := range fields {
		if err := initializers.DB.Create(&field).Error; err != nil {
			log.Fatalf("failed to seed field: %v", err)
		}
	}
}
