package models

import "gorm.io/gorm"

// Notification belongs to User
type Notification struct {
	gorm.Model
	UserID int
	Detail string `gorm:"size:300"`
}
