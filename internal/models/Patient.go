package models

import (
	"gorm.io/gorm"
)

type Patient struct {
	gorm.Model
	First_Name string
	Email      string
}
