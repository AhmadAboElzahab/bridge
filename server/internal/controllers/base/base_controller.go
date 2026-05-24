package base

import (
	"errors"
	"net/http"
	"reflect"

	"github.com/AhmadAboElzahab/bridge/internal/initializers"
	"github.com/AhmadAboElzahab/bridge/internal/services"
	"github.com/AhmadAboElzahab/bridge/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// newModelInstance returns a fresh zero-value pointer of the same type as c.Model,
// avoiding shared-state mutations across concurrent requests.
func (c *BaseController) newModelInstance() interface{} {
	return reflect.New(reflect.TypeOf(c.Model).Elem()).Interface()
}

type BaseController struct {
	Model interface{}
}

func (c *BaseController) Index(ctx *gin.Context) {
	input, err := services.BindTabPayload(ctx)
	if err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid payload", err.Error())
		return
	}

	db := initializers.DB.WithContext(ctx.Request.Context())

	err = db.Transaction(func(tx *gorm.DB) error {
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

		if !field.IsExported() || field.Anonymous {
			continue
		}

		if field.Type.String() == "gorm.DeletedAt" || field.Type.PkgPath() == "gorm.io/gorm" || field.Type.PkgPath() == "time" {
			continue
		}

		kind := field.Type.Kind()

		if kind == reflect.Slice && field.Type.Elem().Kind() == reflect.Struct {
			relations[field.Name] = field.Type
			continue
		}

		if kind == reflect.Struct && field.Type.Name() != "" && field.Type.PkgPath() != "" {
			relations[field.Name] = field.Type
		}
	}

	return relations
}

func (c *BaseController) Store(ctx *gin.Context) {}

func (c *BaseController) Show(ctx *gin.Context) {
	id := ctx.Param("id")
	db := initializers.DB.WithContext(ctx.Request.Context())
	model := c.newModelInstance()
	if err := db.First(model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ErrorJSON(ctx, http.StatusNotFound, "Resource not found")
		} else {
			utils.ErrorJSON(ctx, http.StatusInternalServerError, "Database error", err.Error())
		}
		return
	}
	ctx.JSON(http.StatusOK, model)
}

func (c *BaseController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	db := initializers.DB.WithContext(ctx.Request.Context())
	model := c.newModelInstance()

	if err := db.First(model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ErrorJSON(ctx, http.StatusNotFound, "Resource not found")
		} else {
			utils.ErrorJSON(ctx, http.StatusInternalServerError, "Database error", err.Error())
		}
		return
	}

	var updates map[string]interface{}
	if err := ctx.ShouldBindJSON(&updates); err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	for _, field := range []string{"id", "created_at", "updated_at", "deleted_at", "password", "password_hash"} {
		delete(updates, field)
	}

	if len(updates) == 0 {
		ctx.JSON(http.StatusOK, model)
		return
	}

	if err := db.Model(model).Updates(updates).Error; err != nil {
		utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to update", err.Error())
		return
	}

	if err := db.First(model, id).Error; err != nil {
		utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to reload record", err.Error())
		return
	}

	ctx.JSON(http.StatusOK, model)
}

func (c *BaseController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	db := initializers.DB.WithContext(ctx.Request.Context())
	model := c.newModelInstance()

	if err := db.First(model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ErrorJSON(ctx, http.StatusNotFound, "Resource not found")
		} else {
			utils.ErrorJSON(ctx, http.StatusInternalServerError, "Database error", err.Error())
		}
		return
	}
	if err := db.Delete(model).Error; err != nil {
		utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to delete resource", err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Resource deleted successfully"})
}
