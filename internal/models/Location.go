package models

type Country struct {
	ID   uint   `json:"value" gorm:"primaryKey"`
	Name string `json:"label" gorm:"unique"`
	Code string
}
