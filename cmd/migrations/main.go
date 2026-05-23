package main

import (
	"github.com/AhmadAboElzahab/bridge/internal/initializers"
	"github.com/AhmadAboElzahab/bridge/internal/models"
)

func init() {
	initializers.LoadENV()
	initializers.ConnectDatabase()
}
func main() {
	initializers.DB.AutoMigrate(&models.User{})
	initializers.DB.AutoMigrate(&models.Country{})
	initializers.DB.AutoMigrate(&models.FormField{})
	initializers.DB.AutoMigrate(&models.UserTab{})
	initializers.DB.AutoMigrate(&models.UserTabColumn{})
	initializers.DB.AutoMigrate(&models.Language{})
	initializers.DB.AutoMigrate(&models.Skill{})
	initializers.DB.AutoMigrate(&models.Maid{})
	// plop:models
}
