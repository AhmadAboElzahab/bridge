package routes

import (
	"github.com/AhmadAboElzahab/bridge/controllers/auth"
	"github.com/AhmadAboElzahab/bridge/controllers/patient"
	"github.com/AhmadAboElzahab/bridge/controllers/user"
	"github.com/AhmadAboElzahab/bridge/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Create controllers
	userCtrl := user.NewUserController()
	patientCtrl := patient.NewPatientController()
	authCtrl := auth.NewAuthController()

	api := r.Group("/api")

	{
		{
			auth := api.Group("/auth")
			auth.POST("/signin", authCtrl.Signin)
			auth.POST("/signup", authCtrl.Signup)

		}
		api.Use(middlewares.AuthMiddleware())
		{
			users := api.Group("/users")
			users.GET("/", userCtrl.Index)
			users.POST("/", userCtrl.Store)
			users.GET("/:id", userCtrl.Show)
			users.PUT("/:id", userCtrl.Update)
			users.DELETE("/:id", userCtrl.Delete)
		}

		{
			patients := api.Group("/patients")
			patients.GET("/", patientCtrl.Index)
			patients.POST("/", patientCtrl.Store)
			patients.GET("/:id", patientCtrl.Show)
			patients.PUT("/:id", patientCtrl.Update)
			patients.DELETE("/:id", patientCtrl.Delete)
		}
	}

}
