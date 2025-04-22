package models

type Country struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique"`
	Code string
}
