package base

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"unicode"

	"github.com/AhmadAboElzahab/bridge/internal/initializers"
	"github.com/AhmadAboElzahab/bridge/internal/models"
	"github.com/gin-gonic/gin"
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
		Search  string                 `json:"search"`
		Filters map[string]interface{} `json:"filters" binding:"required"`
		Columns []ColumnInput          `json:"columns" binding:"required,dive"`
		Page    int                    `json:"page"`
		Size    int                    `json:"size"`
	}

	var input TabPayload
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload", "details": err.Error()})
		return
	}
	if input.Page < 0 {
		input.Page = 0
	}
	if input.Size <= 0 || input.Size > 100 {
		input.Size = 50
	}

	modelType := reflect.TypeOf(c.Model).Elem()
	tableName := initializers.DB.NamingStrategy.TableName(modelType.Name())

	// Load tab and save filters
	var tab models.UserTab
	if err := initializers.DB.Where("id = ?", input.TabID).First(&tab).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Tab not found"})
		return
	}
	filtersJSON, _ := json.Marshal(input.Filters)
	tab.Filters = filtersJSON
	initializers.DB.Save(&tab)

	// Load column visibility
	var tabColumns []models.UserTabColumn
	initializers.DB.Preload("FormField").Where("user_tab_id = ?", tab.ID).Find(&tabColumns)
	visibleFields := map[string]bool{}
	for _, col := range tabColumns {
		if col.Visible {
			visibleFields[col.FormField.FieldKey] = true
		}
	}

	// Build query
	query := initializers.DB.Table(tableName).
		Model(c.Model).
		Preload("Skills").
		Preload("Languages").
		Preload("Nationality").
		Preload(clause.Associations)

	// Apply filters
	for key, val := range input.Filters {
		query = query.Where(fmt.Sprintf("%s.%s = ?", tableName, key), val)
	}

	// Relation JOIN map
	relationJoins := map[string]map[string]string{
		"maids": {
			"skills":      "LEFT JOIN maid_skills ON maid_skills.maid_id = maids.id LEFT JOIN skills ON skills.id = maid_skills.skill_id",
			"languages":   "LEFT JOIN maid_languages ON maid_languages.maid_id = maids.id LEFT JOIN languages ON languages.id = maid_languages.language_id",
			"nationality": "LEFT JOIN countries ON countries.id = maids.nationality_id",
		},
	}

	// Dynamic JOIN + text fields
	textFields := getSearchableTextFields(c.Model)
	joined := map[string]bool{}
	if joinMap, ok := relationJoins[tableName]; ok {
		for _, col := range input.Columns {
			if strings.Contains(col.FieldKey, ".") {
				parts := strings.SplitN(col.FieldKey, ".", 2)
				relation, field := parts[0], parts[1]
				if joinStmt, exists := joinMap[relation]; exists && !joined[relation] {
					query = query.Joins(joinStmt)
					joined[relation] = true
				}
				textFields = append(textFields, fmt.Sprintf("%s.%s", relation, field))
			}
		}
		query = query.Group(fmt.Sprintf("%s.id", tableName))
	}

	// Global search
	if input.Search != "" && len(textFields) > 0 {
		var conditions []string
		var args []interface{}
		for _, field := range textFields {
			conditions = append(conditions, fmt.Sprintf("%s ILIKE ?", field))
			args = append(args, "%"+input.Search+"%")
		}
		query = query.Where(strings.Join(conditions, " OR "), args...)
	}

	// Pagination
	var total int64
	query.Count(&total)

	offset := input.Page * input.Size
	query = query.Offset(offset).Limit(input.Size)

	// Execute
	modelSlice := reflect.MakeSlice(reflect.SliceOf(modelType), 0, 0)
	resultPtr := reflect.New(modelSlice.Type()).Interface()
	if err := query.Find(resultPtr).Error; err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "Failed to fetch records"})
		return
	}

	// Convert to generic JSON
	rawJSON, _ := json.Marshal(resultPtr)
	var filtered []map[string]interface{}
	_ = json.Unmarshal(rawJSON, &filtered)

	// Response
	ctx.JSON(http.StatusOK, gin.H{
		"data": filtered,
		"meta": gin.H{
			"totalRowCount": total,
		},
	})
}

func getSearchableTextFields(model any) []string {
	t := reflect.TypeOf(model).Elem()
	tableName := initializers.DB.NamingStrategy.TableName(t.Name())
	var fields []string

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Type.Kind() == reflect.String && f.Tag.Get("json") != "-" {
			column := f.Tag.Get("gorm")
			if column == "" || !strings.Contains(column, "column:") {
				column = f.Tag.Get("json")
			} else {
				for _, part := range strings.Split(column, ";") {
					if strings.HasPrefix(part, "column:") {
						column = strings.TrimPrefix(part, "column:")
						break
					}
				}
			}
			if column == "" {
				column = strings.ToLower(f.Name)
			}
			fields = append(fields, fmt.Sprintf("%s.%s", tableName, column))
		}
	}

	return fields
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
