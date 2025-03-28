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
	router.Static("/storage", "./storage")

	routes.SetupRoutes(router)

	router.Run()
}
