package tabs

import (
	"net/http"
	"sort"

	"github.com/AhmadAboElzahab/bridge/internal/initializers"
	"github.com/AhmadAboElzahab/bridge/internal/models"
	"github.com/AhmadAboElzahab/bridge/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type TabsController struct{}

func NewTabsController() *TabsController {
	return &TabsController{}
}

// GET /api/tabs?model=Patient
func (tc *TabsController) GetTabs(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uint)
	model := ctx.Query("model")

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

		if options, err := utils.ResolveOptionsForField(field); err == nil {
			fieldMap["options"] = options
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
		columns := []gin.H{}
		for _, col := range tab.Columns {
			columns = append(columns, gin.H{
				"form_field_id": col.FormFieldID,
				"field_key":     col.FieldKey,
				"label":         col.FormField.Label,
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
func (tc *TabsController) AddNewTab(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uint)
	type CreateTabInput struct {
		ModelName string `json:"model_name" binding:"required"`
		TabName   string `json:"tab_name" binding:"required"`
	}

	var input CreateTabInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload", "details": err.Error()})
		return
	}

	newTab := models.UserTab{
		UserID:     userID,
		ModelName:  input.ModelName,
		TabName:    input.TabName,
		IsDefault:  false,
		SearchTerm: "",
		Filters:    datatypes.JSON([]byte(`{}`)),
	}
	if err := initializers.DB.Create(&newTab).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user tab"})
		return
	}

	var formFields []models.FormField
	initializers.DB.Where("model_name = ?", input.ModelName).Find(&formFields)
	for i, f := range formFields {
		col := models.UserTabColumn{
			UserTabID:   newTab.ID,
			FormFieldID: &f.ID,
			Visible:     f.TableIsVisible,
			Locked:      f.TableIsPinned,
			Order:       i + 1,
			Width:       f.FormWidth,
		}
		initializers.DB.Create(&col)
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "New tab created", "tab_id": newTab.ID})
}

// PUT /api/tabs/:id
func (tc *TabsController) UpdateTab(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uint)
	tabID := ctx.Param("id")

	var tab models.UserTab
	if err := initializers.DB.Where("id = ? AND user_id = ?", tabID, userID).First(&tab).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Tab not found"})
		return
	}

	var input struct {
		TabName string `json:"tab_name" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload", "details": err.Error()})
		return
	}

	tab.TabName = input.TabName
	if err := initializers.DB.Save(&tab).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tab name"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Tab name updated successfully"})
}
func (tc *TabsController) DeleteTab(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uint)
	tabID := ctx.Param("id")

	var tab models.UserTab
	if err := initializers.DB.Where("id = ? AND user_id = ?", tabID, userID).First(&tab).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Tab not found"})
		return
	}

	initializers.DB.Where("user_tab_id = ?", tab.ID).Delete(&models.UserTabColumn{})
	initializers.DB.Delete(&tab)

	ctx.JSON(http.StatusOK, gin.H{"message": "Tab deleted successfully"})
}
