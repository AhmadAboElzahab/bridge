package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/AhmadAboElzahab/bridge/internal/initializers"
	"github.com/AhmadAboElzahab/bridge/internal/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)


func init() {
	initializers.LoadENV()
	initializers.ConnectDatabase()
	initializers.ConnectStorage()
}

func main() {
	router := gin.Default()

	allowOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if allowOrigins == "" {
		allowOrigins = "http://localhost:3000"
	}
	router.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(allowOrigins, ","),
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	gin.ForceConsoleColor()
	router.Use(gin.Logger())

	// Only serve local files when using the local storage driver.
	if d := strings.ToLower(os.Getenv("STORAGE_DRIVER")); d == "" || d == "local" {
		router.Static("/storage", "./storage")
	}
	routes.SetupRoutes(router)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		log.Println("Server running on http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	if sqlDB, err := initializers.DB.DB(); err == nil {
		sqlDB.Close()
	}

	log.Println("Server stopped")
}
