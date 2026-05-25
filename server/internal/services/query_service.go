package services

import (
	"fmt"
	"reflect"
	"strings"

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

	query = applySorting(query, input.Sorting, formFieldsByKey, modelTable)

	if err := query.Offset((input.Page - 1) * input.Size).Limit(input.Size).Find(slicePtr).Error; err != nil {
		return nil, 0, err
	}

	return slicePtr, total, nil
}

// textSortFieldTypes are field types whose values are free-text strings.
// For these we wrap the column in LOWER() so that sorting is case-insensitive
// and produces the alphabetical order users expect (a < b … z regardless of case).
var textSortFieldTypes = map[string]bool{
	"string_field":             true,
	"email_field":              true,
	"phone_field":              true,
	"url_field":                true,
	"social_url_field":         true,
	"link_field":               true,
	"single_select":            true,
	"radio_select":             true,
	"creatable_single_select":  true,
	"rich_field":               true,
}

// applySorting builds ORDER BY clauses from the client sort state.
// For single_relation fields with a DataSource, it LEFT JOINs the related table
// and sorts by the human-readable label column instead of the FK integer.
// For text fields it wraps the column in LOWER() for case-insensitive ordering.
// Always appends {table}.id ASC as a stable tiebreaker.
func applySorting(
	query *gorm.DB,
	sorting []SortItem,
	formFields map[string]models.FormField,
	modelTable string,
) *gorm.DB {
	joinedKeys := make(map[string]bool)

	for _, s := range sorting {
		field, ok := formFields[s.FieldKey]
		if !ok {
			continue
		}

		dir := "ASC"
		if s.Desc {
			dir = "DESC"
		}

		if field.FormFieldType == "single_relation" && field.DataSource != "" {
			parts := strings.Split(field.DataSource, ":")
			if len(parts) != 3 {
				continue
			}
			joinTable := query.NamingStrategy.TableName(parts[0])
			labelCol := query.NamingStrategy.ColumnName("", parts[2])
			fkCol := field.FieldKey + "_id"

			joinKey := joinTable + "." + fkCol
			if !joinedKeys[joinKey] {
				query = query.Joins(fmt.Sprintf(
					"LEFT JOIN %s ON %s.%s = %s.id",
					joinTable, modelTable, fkCol, joinTable,
				))
				joinedKeys[joinKey] = true
			}
			query = query.Order(fmt.Sprintf("LOWER(%s.%s) %s", joinTable, labelCol, dir))
		} else if textSortFieldTypes[field.FormFieldType] {
			query = query.Order(fmt.Sprintf("LOWER(%s.%s) %s", modelTable, field.FieldKey, dir))
		} else {
			query = query.Order(fmt.Sprintf("%s.%s %s", modelTable, field.FieldKey, dir))
		}
	}

	// Stable tiebreaker so pagination is deterministic
	query = query.Order(modelTable + ".id ASC")
	return query
}
