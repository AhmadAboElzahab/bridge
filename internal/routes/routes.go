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

func SetupRoutes(r *gin.Engine) {
	userCtrl := user.NewUserController()
	patientCtrl := patient.NewPatientController()
	authCtrl := auth.NewAuthController()

	// Swagger docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public API group
	api := r.Group("/api")

	// ✅ Public routes
	authRoutes := api.Group("/auth")
	{
		authRoutes.POST("/signin", authCtrl.Signin)
		authRoutes.POST("/signup", authCtrl.Signup)
	}

	protected := api.Group("/")
	protected.Use(middlewares.AuthMiddleware())

	users := protected.Group("/users")
	{
		users.GET("/", userCtrl.Index)
		users.POST("/", userCtrl.Store)
		users.GET("/:id", userCtrl.Show)
		users.PUT("/:id", userCtrl.Update)
		users.DELETE("/:id", userCtrl.Delete)
	}

	patients := protected.Group("/patients")
	{
		patients.GET("/", patientCtrl.Index)
		patients.POST("/", patientCtrl.Store)
		patients.GET("/:id", patientCtrl.Show)
		patients.PUT("/:id", patientCtrl.Update)
		patients.DELETE("/:id", patientCtrl.Delete)
	}
}
