package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/AhmadAboElzahab/bridge/internal/constants"
	"github.com/AhmadAboElzahab/bridge/internal/initializers"
	"github.com/AhmadAboElzahab/bridge/internal/models"
)

// ResolveOptionsFromDataSource dynamically loads a list of key-label options from a model
// defined by a `FormField`-style DataSource string in the format "Model:ValueField:LabelField".
//
// This is typically used to populate dropdowns/select fields in forms,
// pulling options from related models like countries, skills, or languages.
//
// Supported models: "Country", "Skill", "Language" via constants.
//
// Parameters:
//   - dataSource: a string in the format "Model:IDField:LabelField" (e.g., "Country:ID:Name")
//
// Returns:
//   - a slice of option maps, each with `value` and `label` keys
//   - an error if the model is unsupported or the fields are invalid
func ResolveOptionsFromDataSource(dataSource string) ([]map[string]interface{}, error) {
	parts := strings.Split(dataSource, ":")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid datasource format")
	}

	modelName := parts[0]
	valueField := parts[1]
	labelField := parts[2]
	var result any

	// Select the appropriate model slice based on the model name
	switch modelName {
	case constants.COUNTRIES:
		result = &[]models.Country{}
	case constants.SKILLS:
		result = &[]models.Skill{}
	case constants.LANGUAGES:
		result = &[]models.Language{}
	default:
		return nil, fmt.Errorf("unsupported model for datasource: %s", modelName)
	}

	// Fetch records from the database
	if err := initializers.DB.Find(result).Error; err != nil {
		return nil, err
	}

	// Use reflection to extract the value/label fields from each record
	v := reflect.ValueOf(result).Elem()
	options := []map[string]interface{}{}

	for i := 0; i < v.Len(); i++ {
		item := v.Index(i)
		val := item.FieldByName(valueField)
		label := item.FieldByName(labelField)

		if val.IsValid() && label.IsValid() {
			options = append(options, map[string]interface{}{
				"value": val.Interface(),
				"label": label.Interface(),
			})
		}
	}

	return options, nil
}

// ResolveOptionsForField returns a slice of {"value", "label"} options from either
// the static `Options` JSON or a dynamic `DataSource`.
func ResolveOptionsForField(field models.FormField) ([]map[string]interface{}, error) {
	// If DataSource is provided → resolve dynamically
	if field.DataSource != "" {
		return ResolveOptionsFromDataSource(field.DataSource)
	}

	// If Options is provided → parse JSON
	if len(field.Options) > 0 {
		var opts []map[string]interface{}
		if err := json.Unmarshal(field.Options, &opts); err != nil {
			return nil, fmt.Errorf("invalid options JSON for field '%s': %w", field.FieldKey, err)
		}
		return opts, nil
	}

	// Neither options nor datasource
	return nil, fmt.Errorf("no options or datasource provided for field '%s'", field.FieldKey)
}
