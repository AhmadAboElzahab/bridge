// patient/dto.go
package maid

type PatientInput struct {
	First_Name string `json:"first_name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
}
