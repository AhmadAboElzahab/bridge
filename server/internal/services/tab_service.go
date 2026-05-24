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
func BindTabPayload(ctx *gin.Context) (*TabPayload, error) {
	var input TabPayload
	if err := ctx.ShouldBindJSON(&input); err != nil {
		return nil, err
	}
	if input.Page < 1 {
		input.Page = 1
	}
	if input.Size <= 0 || input.Size > 200 {
		input.Size = 50
	}
	return &input, nil
}

// UpdateTabSettings updates the user tab's filters and search term in the database.
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

// LoadFormFields retrieves all FormField definitions for a given model name.
func LoadFormFields(
	tx *gorm.DB,
	modelName string,
) (map[string]models.FormField, map[string]models.FormField, error) {
	var fields []models.FormField
	if err := tx.Where("model_name = ?", modelName).Find(&fields).Error; err != nil {
		return nil, nil, err
	}

	byFieldKey := make(map[string]models.FormField, len(fields))
	byID := make(map[string]models.FormField, len(fields))

	for _, f := range fields {
		byFieldKey[f.FieldKey] = f
		byID[fmt.Sprintf("%d", f.ID)] = f
	}
	return byFieldKey, byID, nil
}

// UpsertTabColumns creates or updates user tab column settings.
// Uses a single SELECT + batch INSERT/UPDATE strategy instead of N individual queries.
func UpsertTabColumns(tx *gorm.DB, tabID uint, columns []ColumnInput, fieldMap map[string]models.FormField) error {
	if len(columns) == 0 {
		return nil
	}

	// Single query to fetch all existing columns for this tab
	var existing []models.UserTabColumn
	if err := tx.Where("user_tab_id = ?", tabID).Find(&existing).Error; err != nil {
		return err
	}
	existingByFieldID := make(map[uint]*models.UserTabColumn, len(existing))
	for i := range existing {
		if existing[i].FormFieldID != nil {
			existingByFieldID[*existing[i].FormFieldID] = &existing[i]
		}
	}

	var toCreate []models.UserTabColumn
	var toUpdate []models.UserTabColumn

	for _, col := range columns {
		ff, ok := fieldMap[col.FieldKey]
		if !ok {
			return fmt.Errorf("invalid field_key: %s", col.FieldKey)
		}

		if tc, exists := existingByFieldID[ff.ID]; exists {
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
			toUpdate = append(toUpdate, *tc)
		} else {
			tc := models.UserTabColumn{
				UserTabID:   tabID,
				FormFieldID: &ff.ID,
				FieldKey:    col.FieldKey,
			}
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
			toCreate = append(toCreate, tc)
		}
	}

	if len(toCreate) > 0 {
		if err := tx.Create(&toCreate).Error; err != nil {
			return err
		}
	}
	for i := range toUpdate {
		if err := tx.Save(&toUpdate[i]).Error; err != nil {
			return err
		}
	}
	return nil
}
