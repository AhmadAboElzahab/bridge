package models

type Language struct {
	ID       uint   `gorm:"primaryKey"`
	ISO639_3 string `gorm:"size:3;unique;not null"`
	Name     string `gorm:"not null"`
}
