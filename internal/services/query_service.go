package services

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/AhmadAboElzahab/bridge/internal/models"
	"github.com/AhmadAboElzahab/bridge/internal/utils"
	"gorm.io/gorm"
)

// QueryModelRecords dynamically queries records from the given GORM model based on provided filters, search term,
// and pagination input. It supports dynamic relation preloading and flexible search across form fields.
//
// Parameters:
//   - tx: the GORM DB transaction instance
//   - model: a pointer to the model struct type (e.g., &models.Maid{})
//   - input: a TabPayload containing filters, search string, page number, and page size
//   - formFields: a map of form fields used for dynamic search (field_key -> FormField)
//
// Returns:
//   - a slice of maps representing the queried rows (with relation data if preloaded)
//   - the total number of records matching the query (before pagination)
//   - an error if any step fails during the query process
func QueryModelRecords(
	tx *gorm.DB,
	model interface{},
	input *TabPayload,
	formFields map[string]models.FormField,
) ([]map[string]interface{}, int64, error) {

	modelType := reflect.TypeOf(model).Elem()
	modelSlice := reflect.MakeSlice(reflect.SliceOf(modelType), 0, 0)
	resultPtr := reflect.New(modelSlice.Type()).Interface()
	modelTable := tx.NamingStrategy.TableName(modelType.Name())

	query := tx.Model(model).Table(modelTable)

	for rel := range utils.DiscoverRelations(modelType) {
		query = query.Preload(rel)
	}

	// 🔍 Dynamic search
	if input.Search != "" {
		fieldSlice := make([]models.FormField, 0, len(formFields))
		for _, f := range formFields {
			fieldSlice = append(fieldSlice, f)
		}
		sb := utils.NewSearchBuilder(query, modelTable, fieldSlice, input.Search)
		query = sb.Build()
	}

	// 🧪 Filters
	for k, v := range input.Filters {
		query = query.Where(fmt.Sprintf("%s.%s = ?", modelTable, k), v)
	}

	// 📊 Pagination
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Offset(input.Page * input.Size).Limit(input.Size).Find(resultPtr).Error; err != nil {
		return nil, 0, err
	}

	//Convert to JSON then to map[]
	raw, _ := json.Marshal(resultPtr)
	var output []map[string]interface{}
	_ = json.Unmarshal(raw, &output)

	return output, total, nil
}
