package maid

import (
	"net/http"

	"github.com/AhmadAboElzahab/bridge/internal/controllers/base"
	"github.com/AhmadAboElzahab/bridge/internal/initializers"
	"github.com/AhmadAboElzahab/bridge/internal/models"
	"github.com/AhmadAboElzahab/bridge/internal/utils"
	"github.com/gin-gonic/gin"
)

type MaidController struct {
	base.BaseController
}

func NewMaidController() *MaidController {
	return &MaidController{
		BaseController: base.BaseController{
			Model: &models.Maid{},
		},
	}
}

type MaidInput struct {
	FirstName string `json:"first_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
}

// Store godoc
// @Summary Create a new maid
// @Description Create a new maid with first name and email
// @Tags maids
// @Accept json
// @Produce json
// @Param maid body MaidInput true "Maid data"
// @Success 201 {object} models.Maid
// @Failure 400 {object} utils.ErrorResponse
// @Router /maids [post]
func (pc *MaidController) Store(ctx *gin.Context) {
	var input MaidInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: err.Error()})
		return
	}

	maid := models.Maid{
		FirstName: input.FirstName,
		Email:     input.Email,
	}

	if err := initializers.DB.Create(&maid).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: "Failed to create maid"})
		return
	}

	ctx.JSON(http.StatusCreated, maid)
}

// Show godoc
// @Summary Get a maid by ID
// @Description Get a single maid using their ID
// @Tags maids
// @Produce json
// @Param id path int true "Maid ID"
// @Success 200 {object} models.Maid
// @Failure 404 {object} utils.ErrorResponse
// @Router /maids/{id} [get]
func (pc *MaidController) Show(ctx *gin.Context) {
	pc.BaseController.Show(ctx)
}

// Update godoc
// @Summary Update a maid
// @Description Update maid information by ID
// @Tags maids
// @Accept json
// @Produce json
// @Param id path int true "Maid ID"
// @Param maid body MaidInput true "Updated maid data"
// @Success 200 {object} models.Maid
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /maids/{id} [put]
func (pc *MaidController) Update(ctx *gin.Context) {
	pc.BaseController.Update(ctx)
}

// Delete godoc
// @Summary Delete a maid
// @Description Delete a maid by ID
// @Tags maids
// @Param id path int true "Maid ID"
// @Success 200 {object} string
// @Failure 404 {object} utils.ErrorResponse
// @Router /maids/{id} [delete]
func (pc *MaidController) Delete(ctx *gin.Context) {
	pc.BaseController.Delete(ctx)
}
