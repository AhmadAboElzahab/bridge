package routes

import (
	"github.com/AhmadAboElzahab/bridge/controllers/patient"
	"github.com/AhmadAboElzahab/bridge/controllers/user"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Create controllers
	userCtrl := user.NewUserController()
	patientCtrl := patient.NewPatientController()

	// User routes
	r.GET("/users", userCtrl.Index)
	r.POST("/users", userCtrl.Store)
	r.GET("/users/:id", userCtrl.Show)
	r.PUT("/users/:id", userCtrl.Update)
	r.DELETE("/users/:id", userCtrl.Delete)

	// Patient routes
	r.GET("/patients", patientCtrl.Index)
	r.POST("/patients", patientCtrl.Store)
	r.GET("/patients/:id", patientCtrl.Show)
	r.PUT("/patients/:id", patientCtrl.Update)
	r.DELETE("/patients/:id", patientCtrl.Delete)
}
