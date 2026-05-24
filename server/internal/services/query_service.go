package services

import (
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
	formFieldsByKey map[string]models.FormField,
	formFieldsByID map[string]models.FormField,
) (interface{}, int64, error) {
	modelType := reflect.TypeOf(model).Elem()
	modelTable := tx.NamingStrategy.TableName(modelType.Name())
	slicePtr := reflect.New(reflect.SliceOf(modelType)).Interface()

	query := tx.Model(model).Table(modelTable)

	for rel := range utils.DiscoverRelations(modelType) {
		query = query.Preload(rel)
	}

	if input.Search != "" {
		ffSlice := make([]models.FormField, 0, len(formFieldsByKey))
		for _, f := range formFieldsByKey {
			ffSlice = append(ffSlice, f)
		}
		query = utils.NewSearchBuilder(query, modelTable, ffSlice, input.Search).Build()
	}

	if len(input.Filters) > 0 {
		query = ApplyFilters(query, input.Filters, formFieldsByKey)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Offset((input.Page - 1) * input.Size).Limit(input.Size).Find(slicePtr).Error; err != nil {
		return nil, 0, err
	}

	return slicePtr, total, nil
}
