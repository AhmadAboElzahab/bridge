package routes

import (
	_ "github.com/AhmadAboElzahab/bridge/docs"
	"github.com/AhmadAboElzahab/bridge/internal/controllers/auth"
	maid "github.com/AhmadAboElzahab/bridge/internal/controllers/maid"
	"github.com/AhmadAboElzahab/bridge/internal/controllers/tabs"
	"github.com/AhmadAboElzahab/bridge/internal/controllers/user"
	"github.com/AhmadAboElzahab/bridge/internal/middlewares"
	// plop:imports
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(r *gin.Engine) {
	userCtrl := user.NewUserController()
	maidCtrl := maid.NewMaidController()
	authCtrl := auth.NewAuthController()
	tabsCtrl := tabs.NewTabsController()
	// plop:controllers

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
		users.POST("/index", userCtrl.Index)
		users.POST("/", userCtrl.Store)
		users.GET("/:id", userCtrl.Show)
		users.PUT("/:id", userCtrl.Update)
		users.DELETE("/:id", userCtrl.Delete)
	}

	maids := protected.Group("/maids")
	{
		maids.POST("/index", maidCtrl.Index)
		maids.POST("/", maidCtrl.Store)
		maids.GET("/:id", maidCtrl.Show)
		maids.PUT("/:id", maidCtrl.Update)
		maids.DELETE("/:id", maidCtrl.Delete)
	}

	tabsRoutes := protected.Group("/tabs")
	{
		tabsRoutes.GET("", tabsCtrl.GetTabs)
		tabsRoutes.POST("", tabsCtrl.AddNewTab)
		tabsRoutes.PUT("/:id", tabsCtrl.UpdateTab)
		tabsRoutes.DELETE("/:id", tabsCtrl.DeleteTab)
	}
	// plop:routes
}
