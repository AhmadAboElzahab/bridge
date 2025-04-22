package models

import (
	"time"

	"gorm.io/gorm"
)

// Patient represents a patient in the system
//
// @swagger:model Patient
type Patient struct {
	// The ID of the patient
	//
	// required: true
	// example: 1
	ID uint `json:"id"`

	// The time the patient was created
	//
	// example: 2025-04-20T15:04:05Z
	CreatedAt time.Time `json:"created_at"`

	// The time the patient was last updated
	//
	// example: 2025-04-20T15:04:05Z
	UpdatedAt time.Time `json:"updated_at"`

	// The time the patient was deleted
	//
	// swaggerignore: true
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// The patient's first name
	//
	// required: true
	// example: John
	First_Name string `json:"first_name"`

	// The patient's email address
	//
	// required: true
	// example: john@example.com
	Email string `json:"email"`

	// The ID of the patient's nationality (linked to Country model)
	//
	// example: 1
	NationalityID uint `json:"-"` // hide raw FK
	// The country object (related model)
	Nationality Country `json:"nationality" gorm:"foreignKey:NationalityID;references:ID"`
}
