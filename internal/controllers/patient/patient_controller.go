package patient

import (
	"github.com/AhmadAboElzahab/bridge/internal/controllers/base"
	"github.com/AhmadAboElzahab/bridge/internal/initializers"
	"github.com/AhmadAboElzahab/bridge/internal/models"
	"github.com/gin-gonic/gin"
)

type PatientController struct {
	base.BaseController
}

func NewPatientController() *PatientController {
	return &PatientController{
		BaseController: base.BaseController{
			Model: &models.Patient{},
		},
	}
}

func (pc *PatientController) Store(ctx *gin.Context) {
	var body struct {
		First_Name string
		Email      string
	}
	if err := ctx.Bind(&body); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	patient := models.Patient{First_Name: body.First_Name, Email: body.Email}
	result := initializers.DB.Create(&patient)
	if result.Error != nil {
		ctx.JSON(400, gin.H{"error": "Failed to create patient"})
		return
	}
	ctx.JSON(201, patient)
}
