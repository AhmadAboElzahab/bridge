package services

import (
	"encoding/json"
	"fmt"

	"github.com/AhmadAboElzahab/bridge/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ColumnInput represents the structure for updating user-specific column settings.
// Each field is optional and may be updated individually.
type ColumnInput struct {
	FieldKey string `json:"field_key" binding:"required"` // Unique identifier of the form field
	Visible  *bool  `json:"visible"`                      // Visibility state of the column
	Locked   *bool  `json:"locked"`                       // Whether the column is locked in the UI
	Order    *int   `json:"order"`                        // Sort order of the column
	Width    *int   `json:"width"`                        // Width of the column in pixels
}

// TabPayload encapsulates the tab-related state passed from the frontend,
// including filters, column settings, pagination, and search input.
type TabPayload struct {
	TabID   uint            `json:"tab_id" binding:"required"` // ID of the user tab
	Filters json.RawMessage `json:"filters" binding:"required"`
	Columns []ColumnInput   `json:"columns" binding:"required,dive"` // Column configurations for the tab
	Page    int             `json:"page"`                            // Page number for pagination
	Size    int             `json:"size"`                            // Page size for pagination
	Search  string          `json:"search"`                          // Global search term
}

// BindTabPayload binds and validates a TabPayload from the Gin request body.
// It also sets sane defaults for pagination values if missing or invalid.
//
// Returns:
//   - a pointer to the parsed TabPayload
//   - an error if the JSON body is invalid
func BindTabPayload(ctx *gin.Context) (*TabPayload, error) {
	var input TabPayload
	if err := ctx.ShouldBindJSON(&input); err != nil {
		return nil, err
	}
	if input.Page < 0 {
		input.Page = 0
	}
	if input.Size <= 0 || input.Size > 100 {
		input.Size = 50
	}
	return &input, nil
}

// UpdateTabSettings updates the user tab's filters and search term in the database
// based on the input TabPayload.
//
// Parameters:
//   - tx: GORM transaction instance
//   - input: TabPayload with tab ID, filters, and search
//
// Returns:
//   - updated UserTab record
//   - error if tab not found or update fails
func UpdateTabSettings(tx *gorm.DB, input *TabPayload) (*models.UserTab, error) {
	var tab models.UserTab
	if err := tx.First(&tab, input.TabID).Error; err != nil {
		return nil, fmt.Errorf("tab not found: %w", err)
	}
	filtersJSON, err := json.Marshal(input.Filters)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal filters: %w", err)
	}
	tab.Filters = filtersJSON
	tab.SearchTerm = input.Search
	if err := tx.Save(&tab).Error; err != nil {
		return nil, err
	}
	return &tab, nil
}

// LoadFormFields retrieves all FormField definitions for a given model name,
// returning them as a map from field key to FormField object.
//
// Parameters:
//   - tx: GORM DB instance
//   - modelName: name of the target model (e.g., "Maid")
//
// Returns:
//   - map of form field keys to FormField definitions
//   - error if query fails
func LoadFormFields(
	tx *gorm.DB,
	modelName string,
) (map[string]models.FormField, map[string]models.FormField, error) {
	var fields []models.FormField
	if err := tx.Where("model_name = ?", modelName).Find(&fields).Error; err != nil {
		return nil, nil, err
	}

	byFieldKey := make(map[string]models.FormField)
	byID := make(map[string]models.FormField)

	for _, f := range fields {
		byFieldKey[f.FieldKey] = f
		byID[fmt.Sprintf("%d", f.ID)] = f
	}
	return byFieldKey, byID, nil
}

// UpsertTabColumns creates or updates user tab column settings based on input.
// It ensures each column is tied to the correct FormField and UserTab.
//
// Parameters:
//   - tx: GORM DB transaction
//   - tabID: the ID of the user tab being updated
//   - columns: slice of column settings to upsert
//   - fieldMap: map of field keys to FormField definitions
//
// Returns:
//   - error if a field key is invalid or the database operation fails
func UpsertTabColumns(tx *gorm.DB, tabID uint, columns []ColumnInput, fieldMap map[string]models.FormField) error {
	for _, col := range columns {
		ff, ok := fieldMap[col.FieldKey]
		if !ok {
			return fmt.Errorf("invalid field_key: %s", col.FieldKey)
		}
		var tc models.UserTabColumn
		tx.Where("user_tab_id = ? AND form_field_id = ?", tabID, ff.ID).FirstOrInit(&tc)
		tc.UserTabID = tabID
		tc.FormFieldID = &ff.ID
		if col.Visible != nil {
			tc.Visible = *col.Visible
		}
		if col.Locked != nil {
			tc.Locked = *col.Locked
		}
		if col.Order != nil {
			tc.Order = *col.Order
		}
		if col.Width != nil {
			tc.Width = *col.Width
		}
		if err := tx.Save(&tc).Error; err != nil {
			return err
		}
	}
	return nil
}
