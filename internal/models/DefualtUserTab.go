package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type DefaultUserTab struct {
	ID        uint `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	ModelName string `json:"model_name"`
	TabName   string `json:"tab_name"`
	IsDefault bool   `json:"is_default"`

	SearchTerm  string         `json:"search_term"`
	Filters     datatypes.JSON `json:"filters"`
	VisibleCols datatypes.JSON `json:"visible_cols"`
}
