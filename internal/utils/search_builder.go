package utils

import (
	"fmt"
	"strings"

	"github.com/AhmadAboElzahab/bridge/internal/constants"
	"github.com/AhmadAboElzahab/bridge/internal/models"
	"gorm.io/gorm"
)

// SearchBuilder dynamically builds a GORM query that supports multi-field full-text search,
// including support for relation-based fields.
//
// It collects necessary JOINs and WHERE clauses based on FormField metadata and applies
// a search term across eligible fields.
type SearchBuilder struct {
	DB         *gorm.DB           // GORM DB instance used for query building
	BaseTable  string             // Name of the base table (model) being queried
	FormFields []models.FormField // List of form fields to include in the search
	SearchTerm string             // The search term to apply across fields

	Joins  []string      // List of JOIN statements to apply
	Wheres []string      // List of WHERE conditions (OR-combined)
	Args   []interface{} // Query arguments used in the WHERE clause
}

// NewSearchBuilder creates a new SearchBuilder instance.
//
// Parameters:
//   - db: GORM DB instance to apply query logic
//   - baseTable: name of the table being queried (e.g., "maids")
//   - formFields: a slice of FormField metadata to determine searchable fields
//   - searchTerm: the search keyword to look for across fields
//
// Returns:
//   - a pointer to a SearchBuilder instance
func NewSearchBuilder(db *gorm.DB, baseTable string, formFields []models.FormField, searchTerm string) *SearchBuilder {
	return &SearchBuilder{
		DB:         db,
		BaseTable:  baseTable,
		FormFields: formFields,
		SearchTerm: searchTerm,
	}
}

// Build constructs the final GORM query by applying necessary JOINs and WHERE conditions
// derived from the form field definitions and search term.
//
// Returns:
//   - a *gorm.DB object ready to execute with the constructed query logic
func (s *SearchBuilder) Build() *gorm.DB {
	for _, field := range s.FormFields {
		s.processField(field)
	}

	db := s.DB.Table(s.BaseTable)
	for _, join := range s.Joins {
		db = db.Joins(join)
	}
	if len(s.Wheres) > 0 {
		db = db.Where(strings.Join(s.Wheres, " OR "), s.Args...)
	}
	db = db.Distinct()

	return db
}

// processField examines an individual FormField and updates the JOIN and WHERE clauses
// accordingly, depending on the field type.
//
// Supported types:
//   - "string", "email", "rich": performs ILIKE search on the field
//   - "multi_relation", "single_relation": joins related tables and searches their labels
//
// Parameters:
//   - f: the FormField to process
func (s *SearchBuilder) processField(f models.FormField) {
	switch f.FormFieldType {
	case constants.TypeMultiRelation, constants.TypeSingleRelation:
		parts := strings.Split(f.DataSource, ":")
		if len(parts) != 3 {
			return
		}
		relModel, _, relLabel := parts[0], parts[1], parts[2]
		_, relTable, err := GetModelForDataSource(relModel)
		if err != nil {
			return
		}

		// Create a unique alias for each relation
		alias := fmt.Sprintf("%s_%s", relTable, f.FieldKey)

		if f.FormFieldType == constants.TypeMultiRelation {
			baseSingular := strings.TrimSuffix(s.BaseTable, "s")    // maids → maid
			pivot := fmt.Sprintf("%s_%s", baseSingular, f.FieldKey) // maid_skills
			relatedSingular := strings.TrimSuffix(f.FieldKey, "s")  // skills → skill

			s.Joins = append(s.Joins,
				fmt.Sprintf("LEFT JOIN %s ON %s.%s_id = %s.id", pivot, pivot, baseSingular, s.BaseTable))

			s.Joins = append(s.Joins,
				fmt.Sprintf("LEFT JOIN %s AS %s ON %s.id = %s.%s_id", relTable, alias, alias, pivot, relatedSingular))
		} else {
			s.Joins = append(s.Joins,
				fmt.Sprintf("LEFT JOIN %s AS %s ON %s.id = %s.%s_id", relTable, alias, alias, s.BaseTable, f.FieldKey))
		}

		s.Wheres = append(s.Wheres, fmt.Sprintf("%s.%s ILIKE ?", alias, relLabel))
		s.Args = append(s.Args, "%"+s.SearchTerm+"%")

	case constants.TypeStringField, constants.TypeEmailField, constants.TypeRichField:
		s.Wheres = append(s.Wheres, fmt.Sprintf("%s.%s ILIKE ?", s.BaseTable, f.FieldKey))
		s.Args = append(s.Args, "%"+s.SearchTerm+"%")
	}
}
