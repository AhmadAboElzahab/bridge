package base

import (
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/AhmadAboElzahab/bridge/internal/initializers"
	"github.com/AhmadAboElzahab/bridge/internal/services"
	"github.com/AhmadAboElzahab/bridge/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RelationUpdateInput struct {
	FieldKey string `json:"field_key" binding:"required"`
	IDs      []uint `json:"ids"`
}

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

	// Coerce zero-value FK fields to nil so PostgreSQL sets NULL instead of violating
	// the foreign key constraint (0 is Go's uint zero value, not a valid FK target).
	for key, val := range updates {
		if strings.HasSuffix(key, "_id") {
			if num, ok := val.(float64); ok && num == 0 {
				updates[key] = nil
			}
		}
	}

	if len(updates) == 0 {
		ctx.JSON(http.StatusOK, model)
		return
	}

	if err := db.Model(model).Updates(updates).Error; err != nil {
		utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to update", err.Error())
		return
	}

	ctx.JSON(http.StatusOK, model)
}

func (c *BaseController) UpdateRelation(ctx *gin.Context) {
	id := ctx.Param("id")
	db := initializers.DB.WithContext(ctx.Request.Context())

	src := c.newModelInstance()
	if err := db.First(src, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ErrorJSON(ctx, http.StatusNotFound, "Resource not found")
		} else {
			utils.ErrorJSON(ctx, http.StatusInternalServerError, "Database error", err.Error())
		}
		return
	}

	var input RelationUpdateInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	assocName := utils.ToPascalCase(input.FieldKey)

	srcVal := reflect.ValueOf(src).Elem()
	field := srcVal.FieldByName(assocName)
	if !field.IsValid() {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Unknown relation field", assocName)
		return
	}

	elemType := field.Type().Elem()
	// reflect.New gives *[]T — db.Find needs the pointer, Replace needs the slice value
	targetsPtr := reflect.New(reflect.SliceOf(elemType))
	targets := targetsPtr.Interface()

	assoc := db.Model(src).Association(assocName)

	if len(input.IDs) == 0 {
		if err := assoc.Clear(); err != nil {
			utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to clear relations", err.Error())
			return
		}
	} else {
		if err := db.Find(targets, input.IDs).Error; err != nil {
			utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to load related records", err.Error())
			return
		}
		// Dereference so GORM receives []T, not *[]T
		if err := assoc.Replace(targetsPtr.Elem().Interface()); err != nil {
			utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to update relations", err.Error())
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Relations updated"})
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
