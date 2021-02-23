package models

import "gorm.io/gorm"

// Role consists of:
//
// ID 1: QA
//
// ID 2: Developer
type Role struct {
	gorm.Model
	Name string `gorm:"size:50"`
}
