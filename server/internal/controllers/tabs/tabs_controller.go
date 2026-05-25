package tabs

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/AhmadAboElzahab/bridge/internal/constants"
	"github.com/AhmadAboElzahab/bridge/internal/initializers"
	"github.com/AhmadAboElzahab/bridge/internal/models"
	"github.com/AhmadAboElzahab/bridge/internal/services"
	"github.com/AhmadAboElzahab/bridge/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
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

var validModels = map[string]bool{
	constants.MAID: true,
	constants.USER: true,
}

// GetTabs godoc
// @Summary     Get form fields and user tabs
// @Description Returns form field metadata and the authenticated user's tab configurations for a given model
// @Tags        tabs
// @Produce     json
// @Security    BearerAuth
// @Param       model  query   string  true  "Model name (e.g. Maid)"
// @Success     200    {object}  object{form_fields=[]object,tabs=[]object}
// @Failure     400    {object}  utils.ErrorResponse
// @Failure     401    {object}  utils.ErrorResponse
// @Failure     500    {object}  utils.ErrorResponse
// @Router      /api/tabs [get]
func (tc *TabsController) GetTabs(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uint)
	model := ctx.Query("model")

	if !validModels[model] {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid model name")
		return
	}

	db := initializers.DB.WithContext(ctx.Request.Context())

	var formFields []models.FormField
	if err := db.Where("model_name = ?", model).Order("field_order ASC").Find(&formFields).Error; err != nil {
		utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to load form fields", err.Error())
		return
	}

	formFieldsResponse := make([]gin.H, 0, len(formFields))
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
	if err := db.Preload("Columns.FormField").
		Where("user_id = ? AND model_name = ?", userID, model).
		Order("is_default DESC, id ASC").
		Find(&tabs).Error; err != nil {
		utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to load tabs", err.Error())
		return
	}

	if len(tabs) == 0 {
		if err := utils.CreateDefaultTabsForUserModel(db, userID, model); err != nil {
			utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to initialize tabs", err.Error())
			return
		}
		if err := db.Preload("Columns.FormField").
			Where("user_id = ? AND model_name = ?", userID, model).
			Order("is_default DESC, id ASC").
			Find(&tabs).Error; err != nil {
			utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to reload tabs", err.Error())
			return
		}
	}

	tabsResponse := make([]gin.H, 0, len(tabs))
	for _, tab := range tabs {
		columns := make([]gin.H, 0, len(tab.Columns))
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
			"sorting":     tab.Sorting,
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

	if !validModels[input.ModelName] {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid model name")
		return
	}

	db := initializers.DB.WithContext(ctx.Request.Context())

	var tabID uint
	err := db.Transaction(func(tx *gorm.DB) error {
		newTab := models.UserTab{
			UserID:     userID,
			ModelName:  input.ModelName,
			TabName:    input.TabName,
			IsDefault:  false,
			SearchTerm: "",
			Filters:    datatypes.JSON([]byte(`{}`)),
		}
		if err := tx.Create(&newTab).Error; err != nil {
			return err
		}
		tabID = newTab.ID

		var formFields []models.FormField
		if err := tx.Where("model_name = ?", input.ModelName).Find(&formFields).Error; err != nil {
			return err
		}

		columns := make([]models.UserTabColumn, 0, len(formFields))
		for i, f := range formFields {
			fieldID := f.ID
			columns = append(columns, models.UserTabColumn{
				UserTabID:   newTab.ID,
				FormFieldID: &fieldID,
				FieldKey:    f.FieldKey,
				Visible:     f.TableIsVisible,
				Locked:      f.TableIsPinned,
				Order:       i + 1,
				Width:       f.FormWidth,
			})
		}
		if len(columns) > 0 {
			return tx.Create(&columns).Error
		}
		return nil
	})
	if err != nil {
		utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to create tab", err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "New tab created", "tab_id": tabID})
}

// UpdateColumnsInput is the request body for bulk-updating column visibility and order.
type UpdateColumnsInput struct {
	Columns []services.ColumnInput `json:"columns" binding:"required,dive"`
}

// UpdateTabColumns godoc
// @Summary     Update tab column settings
// @Description Updates visibility and order for columns in a user's tab
// @Tags        tabs
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id    path    int                 true  "Tab ID"
// @Param       body  body    UpdateColumnsInput  true  "Column updates"
// @Success     200   {object}  object{message=string}
// @Failure     400   {object}  utils.ErrorResponse
// @Failure     401   {object}  utils.ErrorResponse
// @Failure     404   {object}  utils.ErrorResponse
// @Failure     500   {object}  utils.ErrorResponse
// @Router      /api/tabs/{id}/columns [patch]
func (tc *TabsController) UpdateTabColumns(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uint)
	tabID := ctx.Param("id")

	db := initializers.DB.WithContext(ctx.Request.Context())

	var tab models.UserTab
	if err := db.Where("id = ? AND user_id = ?", tabID, userID).First(&tab).Error; err != nil {
		utils.ErrorJSON(ctx, http.StatusNotFound, "Tab not found")
		return
	}

	var input UpdateColumnsInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid payload", err.Error())
		return
	}

	fieldMap, _, err := services.LoadFormFields(db, tab.ModelName)
	if err != nil {
		utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to load form fields", err.Error())
		return
	}

	if err := services.UpsertTabColumns(db, tab.ID, input.Columns, fieldMap); err != nil {
		utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to update columns", err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Columns updated"})
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

	db := initializers.DB.WithContext(ctx.Request.Context())

	var tab models.UserTab
	if err := db.Where("id = ? AND user_id = ?", tabID, userID).First(&tab).Error; err != nil {
		utils.ErrorJSON(ctx, http.StatusNotFound, "Tab not found")
		return
	}

	var input UpdateTabInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid payload", err.Error())
		return
	}

	tab.TabName = input.TabName
	if err := db.Save(&tab).Error; err != nil {
		utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to update tab", err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Tab updated successfully"})
}

// UpdateTabSortingInput is the request body for persisting sort state.
type UpdateTabSortingInput struct {
	Sorting []map[string]interface{} `json:"sorting"`
}

// UpdateTabSorting godoc
// @Summary     Update tab sort state
// @Description Persists the current column sort state for a user's tab
// @Tags        tabs
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id    path    int                      true  "Tab ID"
// @Param       body  body    UpdateTabSortingInput    true  "Sort state"
// @Success     200   {object}  object{message=string}
// @Failure     400   {object}  utils.ErrorResponse
// @Failure     401   {object}  utils.ErrorResponse
// @Failure     404   {object}  utils.ErrorResponse
// @Failure     500   {object}  utils.ErrorResponse
// @Router      /api/tabs/{id}/sorting [patch]
func (tc *TabsController) UpdateTabSorting(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uint)
	tabID := ctx.Param("id")

	db := initializers.DB.WithContext(ctx.Request.Context())

	var tab models.UserTab
	if err := db.Where("id = ? AND user_id = ?", tabID, userID).First(&tab).Error; err != nil {
		utils.ErrorJSON(ctx, http.StatusNotFound, "Tab not found")
		return
	}

	var input UpdateTabSortingInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid payload", err.Error())
		return
	}

	sortingBytes, err := json.Marshal(input.Sorting)
	if err != nil {
		utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to serialize sorting", err.Error())
		return
	}
	sortingJSON := datatypes.JSON(sortingBytes)

	tab.Sorting = sortingJSON
	if err := db.Save(&tab).Error; err != nil {
		utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to save sorting", err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Sorting updated"})
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

	db := initializers.DB.WithContext(ctx.Request.Context())

	var tab models.UserTab
	if err := db.Where("id = ? AND user_id = ?", tabID, userID).First(&tab).Error; err != nil {
		utils.ErrorJSON(ctx, http.StatusNotFound, "Tab not found")
		return
	}

	// Columns cascade-delete via the OnDelete:CASCADE constraint on UserTabColumn.UserTabID
	if err := db.Delete(&tab).Error; err != nil {
		utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to delete tab", err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Tab deleted successfully"})
}
