package base

import (
	"github.com/AhmadAboElzahab/bridge/initializers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BaseController struct {
	Model interface{}
}

func (c *BaseController) Index(ctx *gin.Context) {
	var result []interface{}
	initializers.DB.Find(&result)
	if len(result) == 0 {
		ctx.JSON(200, gin.H{"message": "no records found"})
		return
	}
	ctx.JSON(200, result)
}

func (c *BaseController) Store(ctx *gin.Context) {
	// Logic to bind data and create the resource (customize based on model)
}

func (c *BaseController) Show(ctx *gin.Context) {
	id := ctx.Param("id")
	var modelInstance interface{}
	result := initializers.DB.First(&modelInstance, id)
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
	// Logic to bind data and update the resource (customize based on model)
}

func (c *BaseController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	var modelInstance interface{}
	if err := initializers.DB.First(&modelInstance, id).Error; err != nil {
		ctx.JSON(404, gin.H{"error": "Resource not found"})
		return
	}
	if err := initializers.DB.Delete(&modelInstance).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to delete resource"})
		return
	}
	ctx.JSON(200, gin.H{"message": "Resource deleted successfully"})
}
