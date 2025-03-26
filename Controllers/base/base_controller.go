package base

import (
	"reflect"

	"github.com/AhmadAboElzahab/bridge/initializers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BaseController struct {
	Model interface{}
}

func (c *BaseController) Index(ctx *gin.Context) {
	// Create a slice of the same type as the model
	modelType := reflect.TypeOf(c.Model).Elem()
	sliceType := reflect.SliceOf(modelType)
	results := reflect.New(sliceType).Elem()

	// Find all records
	if err := initializers.DB.Find(results.Addr().Interface()).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to fetch records"})
		return
	}

	if results.Len() == 0 {
		ctx.JSON(200, gin.H{"message": "No records found"})
		return
	}

	ctx.JSON(200, results.Interface())
}

func (c *BaseController) Store(ctx *gin.Context) {
	// This function can be extended to handle model-specific logic
}

func (c *BaseController) Show(ctx *gin.Context) {
	id := ctx.Param("id")
	// Use type assertion to cast the model to the actual type
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
	// Logic to update model (could be extended to handle dynamic fields)
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
