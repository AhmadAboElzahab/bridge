package models

type Skill struct {
	ID   uint   `json:"value" gorm:"primaryKey"`
	Name string `json:"label" gorm:"unique;not null"`
}
