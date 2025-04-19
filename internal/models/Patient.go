package models

import (
	"time"

	"gorm.io/gorm"
)

// Patient model represents the patient entity in the database
// @Description Patient entity with ID, first name, and email
type Patient struct {
	ID         uint           `json:"id" example:"1"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
	First_Name string         `json:"first_name" example:"John"`
	Email      string         `json:"email" example:"john@example.com"`
}
