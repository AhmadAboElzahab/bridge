package constants

import (
	"encoding/json"
	"log"

	"gorm.io/datatypes"
)

func MustMarshalOptions(options []map[string]string) datatypes.JSON {
	data, err := json.Marshal(options)
	if err != nil {
		log.Fatalf("❌ Failed to marshal options: %v", err)
	}
	return datatypes.JSON(data)
}

var (
	GenderOptions = MustMarshalOptions([]map[string]string{
		{"label": "Male", "value": "male"},
		{"label": "Female", "value": "female"},
		{"label": "Other", "value": "other"},
	})

	WorkPreferencesOptions = MustMarshalOptions([]map[string]string{
		{"label": "Live-in", "value": "live-in"},
		{"label": "Live-out", "value": "live-out"},
		{"label": "Full-time", "value": "full-time"},
		{"label": "Part-time", "value": "part-time"},
	})

	VisaStatusOptions = MustMarshalOptions([]map[string]string{
		{"label": "Canceled", "value": "canceled"},
		{"label": "Other", "value": "other"},
	})

	SubscriptionStatusOptions = MustMarshalOptions([]map[string]string{
		{"label": "Paused", "value": "paused"},
		{"label": "Canceled", "value": "canceled"},
		{"label": "None", "value": "none"},
	})

	StatusOptions = MustMarshalOptions([]map[string]string{
		{"label": "Available", "value": "Available"},
		{"label": "Hired", "value": "hired"},
		{"label": "Canceled", "value": "canceled"},
		{"label": "To be retained", "value": "to be retained"},
	})
	AvailabilityOptions = MustMarshalOptions([]map[string]string{
		{"label": "Immediately", "value": "immediately"},
		{"label": "1 Week", "value": "1_week"},
		{"label": "2 Weeks", "value": "2_weeks"},
		{"label": "1 Month", "value": "1_month"},
	})
	CitiesOptions = MustMarshalOptions([]map[string]string{
		{"label": "Abu Dhabi", "value": "abu_dhabi"},
		{"label": "Dubai", "value": "dubai"},
		{"label": "Sharjah", "value": "sharjah"},
		{"label": "Ajman", "value": "ajman"},
		{"label": "Fujairah", "value": "fujairah"},
		{"label": "Ras Al Khaimah", "value": "ras_al_khaimah"},
		{"label": "Umm Al Quwain", "value": "umm_al_quwain"},
		{"label": "Al Ain", "value": "al_ain"},
		{"label": "Khor Fakkan", "value": "khor_fakkan"},
		{"label": "Kalba", "value": "kalba"},
	})
	SubscriptionPlanOptions = MustMarshalOptions([]map[string]string{
		{"label": "Monthly", "value": "monthly"},
		{"label": "Annual", "value": "annual"},
	})
	ReligionOptions = MustMarshalOptions([]map[string]string{
		{"label": "Islam", "value": "islam"},
		{"label": "Christianity", "value": "christianity"},
		{"label": "Hinduism", "value": "hinduism"},
		{"label": "Buddhism", "value": "buddhism"},
		{"label": "Sikhism", "value": "sikhism"},
		{"label": "Judaism", "value": "judaism"},
		{"label": "Atheism", "value": "atheism"},
		{"label": "Agnostic", "value": "agnostic"},
		{"label": "Other", "value": "other"},
		{"label": "Prefer not to say", "value": "prefer_not_to_say"},
	})
)
