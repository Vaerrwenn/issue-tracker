package models

import (
	"gorm.io/gorm"
)

// User belongs to Role
type User struct {
	gorm.Model
	RoleID   int
	Role     Role
	Name     string `gorm:"size:100"`
	Email    string `gorm:"size:300"`
	Password []byte
}
