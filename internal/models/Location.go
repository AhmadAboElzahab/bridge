package models

type Country struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique"`
	Code string
}

// City model struct for GORM
type City struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"unique"`
	CountryID uint
	Country   Country
}
