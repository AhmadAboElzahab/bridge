package base

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
	"unicode"

	"github.com/AhmadAboElzahab/bridge/internal/initializers"
	"github.com/AhmadAboElzahab/bridge/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BaseController struct {
	Model interface{}
}

func toSnakeCase(str string) string {
	var sb strings.Builder
	for i, r := range str {
		if unicode.IsUpper(r) && i > 0 {
			sb.WriteRune('_')
		}
		sb.WriteRune(unicode.ToLower(r))
	}
	return sb.String()
}

func (c *BaseController) Index(ctx *gin.Context) {
	type ColumnInput struct {
		FieldKey string `json:"field_key" binding:"required"`
		Visible  *bool  `json:"visible"`
		Locked   *bool  `json:"locked"`
		Order    *int   `json:"order"`
		Width    *int   `json:"width"`
	}
	type TabPayload struct {
		TabID   uint                   `json:"tab_id" binding:"required"`
		Filters map[string]interface{} `json:"filters" binding:"required"`
		Columns []ColumnInput          `json:"columns" binding:"required,dive"`
	}

	var input TabPayload
	if err := ctx.ShouldBindJSON(&input); err != nil {
		errMessages := []string{}
		if verrs, ok := err.(validator.ValidationErrors); ok {
			for _, verr := range verrs {
				field := verr.Field()
				switch field {
				case "FieldKey":
					errMessages = append(errMessages, "Each column must include a valid field_key")
				default:
					errMessages = append(errMessages, verr.Namespace()+" is "+verr.Tag())
				}
			}
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload", "details": errMessages})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload", "details": []string{err.Error()}})
		}
		return
	}

	var tab models.UserTab
	if err := initializers.DB.Where("id = ?", input.TabID).First(&tab).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Tab not found"})
		return
	}

	filtersJSON, _ := json.Marshal(input.Filters)
	tab.Filters = filtersJSON
	initializers.DB.Save(&tab)

	var fields []models.FormField
	initializers.DB.Where("model_name = ?", tab.ModelName).Find(&fields)
	fieldKeyMap := map[string]models.FormField{}
	for _, f := range fields {
		fieldKeyMap[strings.ToLower(f.FieldKey)] = f
	}

	inputColumnMap := map[string]ColumnInput{}
	unknownKeys := []string{}
	for _, col := range input.Columns {
		key := strings.ToLower(col.FieldKey)
		inputColumnMap[key] = col
		if _, ok := fieldKeyMap[key]; !ok {
			unknownKeys = append(unknownKeys, col.FieldKey)
		}
	}

	missingKeys := []string{}
	for _, field := range fields {
		if _, ok := inputColumnMap[strings.ToLower(field.FieldKey)]; !ok {
			missingKeys = append(missingKeys, field.FieldKey)
		}
	}

	if len(unknownKeys) > 0 || len(missingKeys) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":          "Invalid column configuration",
			"unknown_fields": unknownKeys,
			"missing_fields": missingKeys,
		})
		return
	}

	for _, field := range fields {
		colInput := inputColumnMap[strings.ToLower(field.FieldKey)]
		visible := false
		locked := false
		order := 0
		width := 0
		if colInput.Visible != nil {
			visible = *colInput.Visible
		}
		if colInput.Locked != nil {
			locked = *colInput.Locked
		}
		if colInput.Order != nil {
			order = *colInput.Order
		}
		if colInput.Width != nil {
			width = *colInput.Width
		}

		var col models.UserTabColumn
		if err := initializers.DB.Where("user_tab_id = ? AND form_field_id = ?", tab.ID, field.ID).First(&col).Error; err == nil {
			col.Visible = visible
			col.Locked = locked
			col.Order = order
			col.Width = width
			initializers.DB.Save(&col)
		} else {
			newCol := models.UserTabColumn{
				UserTabID:   tab.ID,
				FormFieldID: field.ID,
				Visible:     visible,
				Locked:      locked,
				Order:       order,
				Width:       width,
			}
			initializers.DB.Create(&newCol)
		}
	}

	var tabColumns []models.UserTabColumn
	initializers.DB.Preload("FormField").Where("user_tab_id = ?", tab.ID).Find(&tabColumns)
	visibleFields := map[string]bool{}
	for _, col := range tabColumns {
		if col.Visible {
			visibleFields[toSnakeCase(col.FormField.FieldKey)] = true
		}
	}

	modelType := reflect.TypeOf(c.Model).Elem()
	sliceType := reflect.SliceOf(modelType)
	results := reflect.New(sliceType).Elem()

	if err := initializers.DB.Preload(clause.Associations).Find(results.Addr().Interface()).Error; err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "Failed to fetch records"})
		return
	}

	rawJSON, _ := json.Marshal(results.Interface())
	var filtered []map[string]interface{}
	_ = json.Unmarshal(rawJSON, &filtered)

	for i := range filtered {
		for k := range filtered[i] {
			if _, ok := visibleFields[k]; !ok {
				delete(filtered[i], k)
			}
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": filtered,
	})
}

func (c *BaseController) Store(ctx *gin.Context) {}

func (c *BaseController) Show(ctx *gin.Context) {
	id := ctx.Param("id")
	modelInstance := c.Model
	result := initializers.DB.First(modelInstance, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			ctx.JSON(404, gin.H{"error": "Resource not found"})
		} else {
			ctx.JSON(500, gin.H{"error": result.Error.Error()})
		}
		return
	}
	ctx.JSON(200, modelInstance)
}

func (c *BaseController) Update(ctx *gin.Context) {}

func (c *BaseController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	modelInstance := c.Model
	if err := initializers.DB.First(modelInstance, id).Error; err != nil {
		ctx.JSON(404, gin.H{"error": "Resource not found"})
		return
	}
	if err := initializers.DB.Delete(modelInstance).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to delete resource"})
		return
	}
	ctx.JSON(200, gin.H{"message": "Resource deleted successfully"})
}
