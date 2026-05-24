package models

import (
	"gorm.io/datatypes"
)

type UserTab struct {
	ID         uint   `gorm:"primaryKey"`
	UserID     uint   `gorm:"index:idx_user_tab_user_model"`
	ModelName  string `gorm:"index:idx_user_tab_user_model"`
	TabName    string
	IsDefault  bool
	SearchTerm string
	Filters    datatypes.JSON
	Columns    []UserTabColumn `gorm:"foreignKey:UserTabID;constraint:OnDelete:CASCADE"`
}
