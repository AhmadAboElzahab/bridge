package routes

import (
	_ "github.com/AhmadAboElzahab/bridge/docs"
	"github.com/AhmadAboElzahab/bridge/internal/controllers/auth"
	"github.com/AhmadAboElzahab/bridge/internal/controllers/patient"
	"github.com/AhmadAboElzahab/bridge/internal/controllers/user"
	"github.com/AhmadAboElzahab/bridge/internal/middlewares"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Adjust the import path

// gin-swagger middleware
// swagger embed files

func SetupRoutes(r *gin.Engine) {
	userCtrl := user.NewUserController()
	patientCtrl := patient.NewPatientController()
	authCtrl := auth.NewAuthController()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
