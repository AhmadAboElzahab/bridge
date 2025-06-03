package base

import (
	"net/http"
	"reflect"

	"github.com/AhmadAboElzahab/bridge/internal/initializers"
	"github.com/AhmadAboElzahab/bridge/internal/services"
	"github.com/AhmadAboElzahab/bridge/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BaseController struct {
	Model interface{}
}

func (c *BaseController) Index(ctx *gin.Context) {
	input, err := services.BindTabPayload(ctx)
	if err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid payload", err.Error())
		return
	}

	err = initializers.DB.Transaction(func(tx *gorm.DB) error {
		tab, err := services.UpdateTabSettings(tx, input)
		if err != nil {
			return err
		}

		fieldMap, fieldMapByID, err := services.LoadFormFields(tx, tab.ModelName)
		if err != nil {
			return err
		}

		if err := services.UpsertTabColumns(tx, tab.ID, input.Columns, fieldMap); err != nil {
			return err
		}

		output, total, err := services.QueryModelRecords(tx, c.Model, input, fieldMap, fieldMapByID)
		if err != nil {
			return err
		}

		ctx.JSON(http.StatusOK, gin.H{
			"data": output,
			"meta": gin.H{"totalRowCount": total},
		})
		return nil
	})

	if err != nil {
		utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to load data", err.Error())
	}
}

func discoverRelations(t reflect.Type) map[string]reflect.Type {
	relations := map[string]reflect.Type{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Skip unexported fields or embedded structs
		if !field.IsExported() || field.Anonymous {
			continue
		}

		// Skip GORM internal and time fields
		if field.Type.String() == "gorm.DeletedAt" || field.Type.PkgPath() == "gorm.io/gorm" || field.Type.PkgPath() == "time" {
			continue
		}

		kind := field.Type.Kind()

		// Handle slices: many-to-many or one-to-many
		if kind == reflect.Slice && field.Type.Elem().Kind() == reflect.Struct {
			relations[field.Name] = field.Type
			continue
		}

		// Handle structs: one-to-one or many-to-one
		if kind == reflect.Struct && field.Type.Name() != "" && field.Type.PkgPath() != "" {
			relations[field.Name] = field.Type
		}
	}

	return relations
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
