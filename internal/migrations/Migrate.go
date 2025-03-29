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
	initializers.DB.AutoMigrate(&models.Patient{})
	initializers.DB.AutoMigrate(&models.User{})

}
