package models

import (
	"gorm.io/datatypes"
)

type UserTab struct {
	ID         uint `gorm:"primaryKey"`
	UserID     uint `gorm:"index"`
	ModelName  string
	TabName    string
	IsDefault  bool
	SearchTerm string
	Filters    datatypes.JSON
	Columns    []UserTabColumn `gorm:"foreignKey:UserTabID;constraint:OnDelete:CASCADE"`
}
