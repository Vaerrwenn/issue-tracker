package models

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	name string `gorm:"size:50"`
}
