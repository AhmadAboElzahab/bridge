package models

import (
	"time"

	"gorm.io/gorm"
)

// User model represents the user entity in the database
// @Description User entity with ID, first name, last name, email, password, date of birth, avatar, and blurhash
type User struct {
	ID          uint           `json:"id" example:"1"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
	FirstName   string         `json:"first_name" example:"John"`
	LastName    string         `json:"last_name" example:"Doe"`
	Email       string         `json:"email" example:"john@example.com"`
	Password    string         `json:"-" swaggerignore:"true"`
	DateOfBirth string         `json:"date_of_birth" example:"1990-01-01"`
	Avatar      string         `json:"avatar" example:"/images/avatar.png"`
	Blurhash    string         `json:"blurhash" example:"LKO2?V%2Tw=w]~RBVZRi};RPxuwH"`
}
