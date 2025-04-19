// patient/dto.go
package patient

type PatientInput struct {
	First_Name string `json:"first_name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
}
