package main

import (
	"fmt"

	"github.com/AhmadAboElzahab/bridge/internal/initializers"
	"github.com/AhmadAboElzahab/bridge/internal/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadENV()
	initializers.ConnectDatabase()
}

func main() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))
	gin.ForceConsoleColor()

	router.Use(gin.Logger())

	router.Static("/storage", "./storage")

	// ✅ Setup all routes
	routes.SetupRoutes(router)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	router.Use(func(c *gin.Context) {
		fmt.Println("🔥", c.Request.Method, c.Request.URL.Path)
		c.Next()
	})

	fmt.Println("Server running on http://localhost:8080")
	router.Run(":8080")
}
