package model

import (
	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	Name     string `gorm:"uniqueIndex:idx_user_file"`
	Visible  bool   `gorm:"default:false"`
	Username string `gorm:"uniqueIndex:idx_user_file"`
}
