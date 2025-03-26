package main

import (
	"github.com/AhmadAboElzahab/bridge/initializers"
	"github.com/AhmadAboElzahab/bridge/models"
)

func init() {
	initializers.LoadENV()
	initializers.ConnectDatabase()
}
func main() {
	initializers.DB.AutoMigrate(&models.Patient{})
	initializers.DB.AutoMigrate(&models.User{})

}
