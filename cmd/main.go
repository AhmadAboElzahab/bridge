package main

import (
	"github.com/AhmadAboElzahab/bridge/internal/initializers"
	"github.com/AhmadAboElzahab/bridge/internal/routes"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadENV()
	initializers.ConnectDatabase()
}
func main() {
	router := gin.Default()
	gin.ForceConsoleColor()
	router.Use(gin.Logger())

	router.Static("/storage", "../../storage")

	routes.SetupRoutes(router)

	router.Run()
}
