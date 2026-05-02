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

// CreateTabInput is the request body for creating a new tab.
type CreateTabInput struct {
	ModelName string `json:"model_name" binding:"required"`
	TabName   string `json:"tab_name" binding:"required"`
}

// UpdateTabInput is the request body for renaming a tab.
type UpdateTabInput struct {
	TabName string `json:"tab_name" binding:"required"`
}

// GetTabs godoc
// @Summary     Get form fields and user tabs
// @Description Returns form field metadata and the authenticated user's tab configurations for a given model
// @Tags        tabs
// @Produce     json
// @Security    BearerAuth
// @Param       model  query   string  true  "Model name (e.g. Maid)"
// @Success     200    {object}  object{form_fields=[]object,tabs=[]object}
// @Failure     401    {object}  utils.ErrorResponse
// @Failure     500    {object}  utils.ErrorResponse
// @Router      /api/tabs [get]
func (tc *TabsController) GetTabs(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uint)
	model := ctx.Query("model")

	var formFields []models.FormField
	if err := initializers.DB.
		Where("model_name = ?", model).
		Order("field_order ASC").
		Find(&formFields).Error; err != nil {
		utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to load form fields", err.Error())
		return
	}

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

// AddNewTab godoc
// @Summary     Create a new tab
// @Description Creates a new table view tab for the authenticated user
// @Tags        tabs
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       body  body    CreateTabInput  true  "Tab creation data"
// @Success     201   {object}  object{message=string,tab_id=int}
// @Failure     400   {object}  utils.ErrorResponse
// @Failure     401   {object}  utils.ErrorResponse
// @Failure     500   {object}  utils.ErrorResponse
// @Router      /api/tabs [post]
func (tc *TabsController) AddNewTab(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uint)

	var input CreateTabInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid payload", err.Error())
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
		utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to create tab", err.Error())
		return
	}

	var formFields []models.FormField
	initializers.DB.Where("model_name = ?", input.ModelName).Find(&formFields)
	for i, f := range formFields {
		col := models.UserTabColumn{
			UserTabID:   newTab.ID,
			FormFieldID: &f.ID,
			FieldKey:    f.FieldKey,
			Visible:     f.TableIsVisible,
			Locked:      f.TableIsPinned,
			Order:       i + 1,
			Width:       f.FormWidth,
		}
		initializers.DB.Create(&col)
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "New tab created", "tab_id": newTab.ID})
}

// UpdateTab godoc
// @Summary     Rename a tab
// @Description Updates the name of a user's tab
// @Tags        tabs
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id    path    int            true  "Tab ID"
// @Param       body  body    UpdateTabInput true  "New tab name"
// @Success     200   {object}  object{message=string}
// @Failure     400   {object}  utils.ErrorResponse
// @Failure     401   {object}  utils.ErrorResponse
// @Failure     404   {object}  utils.ErrorResponse
// @Failure     500   {object}  utils.ErrorResponse
// @Router      /api/tabs/{id} [put]
func (tc *TabsController) UpdateTab(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uint)
	tabID := ctx.Param("id")

	var tab models.UserTab
	if err := initializers.DB.Where("id = ? AND user_id = ?", tabID, userID).First(&tab).Error; err != nil {
		utils.ErrorJSON(ctx, http.StatusNotFound, "Tab not found")
		return
	}

	var input UpdateTabInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid payload", err.Error())
		return
	}

	tab.TabName = input.TabName
	if err := initializers.DB.Save(&tab).Error; err != nil {
		utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to update tab", err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Tab updated successfully"})
}

// DeleteTab godoc
// @Summary     Delete a tab
// @Description Deletes a user's tab and all its column configurations
// @Tags        tabs
// @Produce     json
// @Security    BearerAuth
// @Param       id    path    int  true  "Tab ID"
// @Success     200   {object}  object{message=string}
// @Failure     401   {object}  utils.ErrorResponse
// @Failure     404   {object}  utils.ErrorResponse
// @Failure     500   {object}  utils.ErrorResponse
// @Router      /api/tabs/{id} [delete]
func (tc *TabsController) DeleteTab(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uint)
	tabID := ctx.Param("id")

	var tab models.UserTab
	if err := initializers.DB.Where("id = ? AND user_id = ?", tabID, userID).First(&tab).Error; err != nil {
		utils.ErrorJSON(ctx, http.StatusNotFound, "Tab not found")
		return
	}

	// Columns cascade-delete via the OnDelete:CASCADE constraint on UserTabColumn.UserTabID
	if err := initializers.DB.Delete(&tab).Error; err != nil {
		utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to delete tab", err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Tab deleted successfully"})
}
