package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type UserTab struct {
	ID        uint `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	UserID    uint   `json:"user_id" gorm:"index"`
	ModelName string `json:"model_name"` // e.g., "Patient"
	TabName   string `json:"tab_name"`   // e.g., "Active Patients"
	IsDefault bool   `json:"is_default"` // True for default tab per model per user

	SearchTerm  string         `json:"search_term"`
	Filters     datatypes.JSON `json:"filters"`      // JSON object: map[string]interface{}
	VisibleCols datatypes.JSON `json:"visible_cols"` // JSON array: []string
}
