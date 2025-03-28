package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	First_Name    string
	Last_Name     string
	Email         string
	Password      string
	Date_of_Birth string
	Avatar        string
	Blurhash      string
}
