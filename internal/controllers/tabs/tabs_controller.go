package tabs

import (
	"net/http"

	"github.com/AhmadAboElzahab/bridge/internal/initializers"
	"github.com/AhmadAboElzahab/bridge/internal/models"
	"github.com/AhmadAboElzahab/bridge/internal/utils"
	"github.com/gin-gonic/gin"
)

type TabsController struct{}

func NewTabsController() *TabsController {
	return &TabsController{}
}

// GET /api/tabs?model=Patient
func (tc *TabsController) GetTabs(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uint)
	model := ctx.Query("model")

	var tabs []models.UserTab
	result := initializers.DB.
		Where("user_id = ? AND model_name = ?", userID, model).
		Order("is_default DESC, id ASC").
		Find(&tabs)

	if result.RowsAffected == 0 {
		utils.CreateDefaultTabsForUserModel(userID, model)

		initializers.DB.
			Where("user_id = ? AND model_name = ?", userID, model).
			Order("is_default DESC, id ASC").
			Find(&tabs)
	}

	ctx.JSON(http.StatusOK, tabs)
}

// PUT /api/tabs/:id
func (tc *TabsController) UpdateTab(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uint)
	tabID := ctx.Param("id")

	var tab models.UserTab
	if err := initializers.DB.Where("id = ? AND user_id = ?", tabID, userID).First(&tab).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Tab not found"})
		return
	}

	var input models.UserTab
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tab.SearchTerm = input.SearchTerm
	tab.Filters = input.Filters
	tab.VisibleCols = input.VisibleCols

	initializers.DB.Save(&tab)
	ctx.JSON(http.StatusOK, tab)
}

// DELETE /api/tabs/:id
func (tc *TabsController) DeleteTab(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uint)
	tabID := ctx.Param("id")

	var tab models.UserTab
	if err := initializers.DB.Where("id = ? AND user_id = ?", tabID, userID).First(&tab).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Tab not found"})
		return
	}

	initializers.DB.Delete(&tab)
	ctx.JSON(http.StatusOK, gin.H{"message": "Tab deleted"})
}
