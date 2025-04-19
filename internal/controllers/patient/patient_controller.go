package patient

import (
	"net/http"

	"github.com/AhmadAboElzahab/bridge/internal/controllers/base"
	"github.com/AhmadAboElzahab/bridge/internal/initializers"
	"github.com/AhmadAboElzahab/bridge/internal/models"
	"github.com/AhmadAboElzahab/bridge/internal/utils"
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

// Store godoc
// @Summary Create a new patient
// @Description Create a new patient with first name and email
// @Tags patients
// @Accept json
// @Produce json
// @Param patient body PatientInput true "Patient data"
// @Success 201 {object} models.Patient
// @Failure 400 {object} utils.ErrorResponse
// @Router /patients [post]
func (pc *PatientController) Store(ctx *gin.Context) {
	var body struct {
		First_Name string
		Email      string
	}
	if err := ctx.Bind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: err.Error()})
		return
	}

	patient := models.Patient{First_Name: body.First_Name, Email: body.Email}
	result := initializers.DB.Create(&patient)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: "Failed to create patient"})
		return
	}
	ctx.JSON(201, patient)
}

// Show godoc (inherited from BaseController)
// @Summary Get a patient by ID
// @Description Get a single patient using their ID
// @Tags patients
// @Produce json
// @Param id path int true "Patient ID"
// @Success 200 {object} models.Patient
// @Failure 404 {object} utils.ErrorResponse
// @Router /patients/{id} [get]
func (pc *PatientController) Show(ctx *gin.Context) {
	pc.BaseController.Show(ctx)
}

// Update godoc (inherited from BaseController)
// @Summary Update a patient
// @Description Update patient information by ID
// @Tags patients
// @Accept json
// @Produce json
// @Param id path int true "Patient ID"
// @Param patient body PatientInput true "Updated patient data"
// @Success 200 {object} models.Patient
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /patients/{id} [put]
func (pc *PatientController) Update(ctx *gin.Context) {
	pc.BaseController.Update(ctx)
}

// Delete godoc (inherited from BaseController)
// @Summary Delete a patient
// @Description Delete a patient by ID
// @Tags patients
// @Param id path int true "Patient ID"
// @Success 200 {object} string
// @Failure 404 {object} utils.ErrorResponse
// @Router /api/patients/{id} [delete]
func (pc *PatientController) Delete(ctx *gin.Context) {
	pc.BaseController.Delete(ctx)
}
