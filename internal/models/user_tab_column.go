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

	UserTabID   uint  `gorm:"uniqueIndex:idx_user_tab_column_unique;not null"`
	FormFieldID *uint `gorm:"uniqueIndex:idx_user_tab_column_unique"`
	FieldKey    string `gorm:"column:field_key"`
	Visible     bool
	Locked      bool
	Order       int
	Width       int

	FormField *FormField `gorm:"foreignKey:FormFieldID"`
}
