package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type FormField struct {
	ID        uint `json:"id" gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Core attributes
	Label      string         `json:"label"`
	FieldKey   string         `json:"field_key"`
	ModelName  string         `json:"model_name"`
	HelpText   string         `json:"form_help_text"`
	DataSource string         `json:"data_source"`
	Options    datatypes.JSON `json:"options"`

	// Form display attributes
	FormFieldType  string `json:"form_field_type"`
	FormIsRequired bool   `json:"form_is_required"`
	FormOrder      int    `json:"form_order" gorm:"column:field_order"`
	FormStage      string `json:"form_stage"`
	FormWidth      int    `json:"form_width"`

	TableIsVisible bool `json:"table_is_visible"`
	TableOrder     int  `json:"table_order"`
	TableIsPinned  bool `json:"table_is_pinned"`
}
