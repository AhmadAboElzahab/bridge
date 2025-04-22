package base

import (
	"net/http"
	"reflect"

	"github.com/AhmadAboElzahab/bridge/internal/initializers"
	"github.com/AhmadAboElzahab/bridge/internal/utils"

	"github.com/AhmadAboElzahab/bridge/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BaseController struct {
	Model interface{}
}

func (c *BaseController) Index(ctx *gin.Context) {
	modelType := reflect.TypeOf(c.Model).Elem()
	sliceType := reflect.SliceOf(modelType)
	results := reflect.New(sliceType).Elem()

	// Fetch actual data from database
	if err := initializers.DB.Preload(clause.Associations).Find(results.Addr().Interface()).Error; err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "Failed to fetch records"})
		return
	}

	modelName := modelType.Name()
	var fields []models.FormField
	if err := initializers.DB.
		Where("model_name = ?", modelName).
		Order("field_order ASC").
		Find(&fields).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch form fields",
			"details": err.Error(),
		})
		return
	}

	// Dynamically attach options from DataSource if present
	formFieldsWithOptions := []gin.H{}
	for _, field := range fields {
		fieldMap := gin.H{
			"id":               field.ID,
			"label":            field.Label,
			"field_key":        field.FieldKey,
			"form_field_type":  field.FormFieldType,
			"data_source":      field.DataSource,
			"form_width":       field.FormWidth,
			"form_order":       field.FormOrder,
			"form_stage":       field.FormStage,
			"form_is_required": field.FormIsRequired,
			"table_is_pinned":  field.TableIsPinned,
			"table_is_visible": field.TableIsVisible,
			"table_order":      field.TableOrder,
			"help_text":        field.HelpText,
		}

		// Handle datasource if available
		if field.DataSource != "" {
			options, err := utils.ResolveOptionsFromDataSource(field.DataSource)
			if err == nil {
				fieldMap["options"] = options
			}
		}

		formFieldsWithOptions = append(formFieldsWithOptions, fieldMap)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":        results.Interface(),
		"form_fields": formFieldsWithOptions,
	})
}

func (c *BaseController) Store(ctx *gin.Context) {
}

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

func (c *BaseController) Update(ctx *gin.Context) {
}

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
