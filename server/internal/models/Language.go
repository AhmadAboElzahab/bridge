package models

type Language struct {
	ID       uint   `json:"value" gorm:"primaryKey"`
	ISO639_3 string `gorm:"size:3;unique;not null"`
	Name     string `json:"label" gorm:"not null"`
}
