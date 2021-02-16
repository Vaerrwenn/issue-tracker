package models

import "gorm.io/gorm"

// Notification belongs to User
type Notification struct {
	gorm.Model
	UserID int
	User   User
	Detail string `gorm:"size:300"`
}
