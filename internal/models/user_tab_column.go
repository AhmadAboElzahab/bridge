package models

import (
	"time"

	"gorm.io/gorm"
)

type UserTabColumn struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	UserTabID   uint `gorm:"index"`
	FormFieldID uint `gorm:"index"`
	Visible     bool
	Locked      bool
	Order       int
	Width       int

	FormField FormField `gorm:"foreignKey:FormFieldID"`
}
