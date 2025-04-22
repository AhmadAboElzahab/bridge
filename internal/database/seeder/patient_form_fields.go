package seeder

import (
	"log"

	"github.com/AhmadAboElzahab/bridge/internal/constants"
	"github.com/AhmadAboElzahab/bridge/internal/initializers"
	"github.com/AhmadAboElzahab/bridge/internal/models"
)

func SeedPatientFormFields() {
	fields := []models.FormField{
		{
			Label:     "First Name",
			FieldKey:  "FirstName",
			ModelName: constants.PATIENT,
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
			ModelName: constants.PATIENT,
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
			Label:          "Nationality",
			FieldKey:       "Nationality",
			ModelName:      constants.PATIENT,
			HelpText:       "Select a country",
			FormFieldType:  constants.TypeSingleSelect,
			DataSource:     "Country:ID:Name",
			FormWidth:      2,
			FormOrder:      3,
			FormStage:      "info",
			FormIsRequired: false,
			TableIsPinned:  false,
			TableIsVisible: true,
			TableOrder:     3,
		},
	}

	for _, field := range fields {
		if err := initializers.DB.Create(&field).Error; err != nil {
			log.Fatalf("failed to seed field: %v", err)
		}
	}
}
