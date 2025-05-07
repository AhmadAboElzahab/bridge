package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/AhmadAboElzahab/bridge/internal/constants"
	"github.com/AhmadAboElzahab/bridge/internal/initializers"
	"github.com/AhmadAboElzahab/bridge/internal/models"
)

func ResolveOptionsFromDataSource(dataSource string) ([]map[string]interface{}, error) {
	parts := strings.Split(dataSource, ":")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid datasource format")
	}

	modelName := parts[0]
	valueField := parts[1]
	labelField := parts[2]
	var result any
	switch modelName {
	case constants.COUNTRIES:
		result = &[]models.Country{}

	default:
		return nil, fmt.Errorf("unsupported model for datasource: %s", modelName)
	}

	if err := initializers.DB.Find(result).Error; err != nil {
		return nil, err
	}

	// Reflect over result and extract label/value fields
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
