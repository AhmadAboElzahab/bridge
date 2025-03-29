package user

import (
	"github.com/AhmadAboElzahab/bridge/internal/controllers/base"
	"github.com/AhmadAboElzahab/bridge/internal/initializers"
	"github.com/AhmadAboElzahab/bridge/internal/models"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	base.BaseController
}

func NewUserController() *UserController {
	return &UserController{
		BaseController: base.BaseController{
			Model: &models.User{},
		},
	}
}

func (uc *UserController) Store(ctx *gin.Context) {
	var body struct {
		First_Name string
		Email      string
	}
	if err := ctx.Bind(&body); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	user := models.User{First_Name: body.First_Name, Email: body.Email}
	result := initializers.DB.Create(&user)
	if result.Error != nil {
		ctx.JSON(400, gin.H{"error": "Failed to create user"})
		return
	}
	ctx.JSON(201, user)
}
