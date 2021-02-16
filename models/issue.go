package models

import "gorm.io/gorm"

// Issue Belongs To User
type Issue struct {
	gorm.Model
	UserID   int
	User     User
	Title    string `gorm:"size:100"`
	Body     string `gorm:"size:2000"`
	Status   string `gorm:"size:1"` // 1 = Opened or 0 = Closed
	Severity string `gorm:"size:1"` // 1 = Low, 2 = Medium, 3 = High
}
