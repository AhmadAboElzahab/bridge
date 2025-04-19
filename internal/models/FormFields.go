package models

import (
	"time"

	"gorm.io/gorm"
)

type FormField struct {
	ID        uint `json:"id" gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Label            string            `json:"label"`
	FieldKey         string            `json:"field_key"`
	ModelName        string            `json:"model_name"`  // e.g., "User", "Patient"
	FieldType        string            `json:"field_type"`  // e.g., "string", "date"
	WidgetType       string            `json:"widget_type"` // e.g., "text", "email", "date", "dropdown"
	IsFilterable     bool              `json:"is_filterable"`
	IsRequired       bool              `json:"is_required"`
	Order            int               `json:"order" gorm:"column:field_order"`
	FormFieldOptions []FormFieldOption `json:"options" gorm:"foreignKey:FormFieldID"`
}

type FormFieldOption struct {
	ID        uint `json:"id" gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	FormFieldID uint   `json:"form_field_id"`
	Label       string `json:"label"`
	Value       string `json:"value"`

	// Optional relation
	FormField FormField `gorm:"foreignKey:FormFieldID"`
}
