package tabs

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/AhmadAboElzahab/bridge/internal/initializers"
	"github.com/AhmadAboElzahab/bridge/internal/models"
	"github.com/AhmadAboElzahab/bridge/internal/utils"
	"github.com/gin-gonic/gin"
)

type TabsController struct{}

func NewTabsController() *TabsController {
	return &TabsController{}
}

// GET /api/tabs?model=Patient
func (tc *TabsController) GetTabs(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uint)
	model := ctx.Query("model")

	fmt.Println("Loading tabs for user", userID, "model", model)

	// Load shared form field metadata
	var formFields []models.FormField
	if err := initializers.DB.
		Where("model_name = ?", model).
		Order("field_order ASC").
		Find(&formFields).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load form fields"})
		return
	}

	// Resolve options for fields with data sources
	formFieldsResponse := []gin.H{}
	for _, field := range formFields {
		fieldMap := gin.H{
			"id":               field.ID,
			"label":            field.Label,
			"field_key":        field.FieldKey,
			"model_name":       field.ModelName,
			"form_help_text":   field.HelpText,
			"data_source":      field.DataSource,
			"form_field_type":  field.FormFieldType,
			"form_is_required": field.FormIsRequired,
			"form_order":       field.FormOrder,
			"form_stage":       field.FormStage,
			"form_width":       field.FormWidth,
			"table_is_visible": field.TableIsVisible,
			"table_order":      field.TableOrder,
			"table_is_pinned":  field.TableIsPinned,
		}

		if field.DataSource != "" {
			if options, err := utils.ResolveOptionsFromDataSource(field.DataSource); err == nil {
				fieldMap["options"] = options
			}
		}

		formFieldsResponse = append(formFieldsResponse, fieldMap)
	}

	// Load user tabs with column settings
	var tabs []models.UserTab
	result := initializers.DB.
		Preload("Columns.FormField").
		Where("user_id = ? AND model_name = ?", userID, model).
		Order("is_default DESC, id ASC").
		Find(&tabs)
	fmt.Println("Loaded tabs:", len(tabs))

	if result.RowsAffected == 0 {
		utils.CreateDefaultTabsForUserModel(userID, model)
		initializers.DB.
			Preload("Columns.FormField").
			Where("user_id = ? AND model_name = ?", userID, model).
			Order("is_default DESC, id ASC").
			Find(&tabs)
	}

	tabsResponse := []gin.H{}
	for _, tab := range tabs {
		fmt.Printf("Tab ID: %d has %d columns\n", tab.ID, len(tab.Columns))
		columns := []gin.H{}
		for _, col := range tab.Columns {
			columns = append(columns, gin.H{
				"form_field_id": col.FormFieldID,
				"field_key":     col.FormField.FieldKey,
				"visible":       col.Visible,
				"locked":        col.Locked,
				"order":         col.Order,
				"width":         col.Width,
			})
		}
		sort.Slice(columns, func(i, j int) bool {
			return columns[i]["order"].(int) < columns[j]["order"].(int)
		})

		tabsResponse = append(tabsResponse, gin.H{
			"id":          tab.ID,
			"tab_name":    tab.TabName,
			"model_name":  tab.ModelName,
			"search_term": tab.SearchTerm,
			"filters":     tab.Filters,
			"is_default":  tab.IsDefault,
			"columns":     columns,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"form_fields": formFieldsResponse,
		"tabs":        tabsResponse,
	})
}
