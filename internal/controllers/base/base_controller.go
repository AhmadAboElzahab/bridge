package base

import (
	"net/http"
	"reflect"

	"github.com/AhmadAboElzahab/bridge/internal/initializers"
	"github.com/AhmadAboElzahab/bridge/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BaseController struct {
	Model interface{}
}

func (c *BaseController) Index(ctx *gin.Context) {
	// Get the model type (like Patient or User)
	modelType := reflect.TypeOf(c.Model).Elem()
	sliceType := reflect.SliceOf(modelType)
	results := reflect.New(sliceType).Elem()

	// Fetch actual data from database
	if err := initializers.DB.Find(results.Addr().Interface()).Error; err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "Failed to fetch records"})
		return
	}

	// Fetch form fields for this model
	modelName := modelType.Name()

	var fields []models.FormField
	if err := initializers.DB.
		Preload("FormFieldOptions").
		Where("model_name = ?", modelName).
		Order("field_order ASC").
		Find(&fields).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch form fields",
			"details": err.Error(), // <--- this shows the DB error
		})
		return
	}

	// Return both
	ctx.JSON(http.StatusOK, gin.H{
		"data":        results.Interface(),
		"form_fields": fields,
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
