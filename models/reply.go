package models

import (
	"issue-tracker/database"

	"gorm.io/gorm"
)

// Reply on each Issue and Users.
type Reply struct {
	gorm.Model
	UserID  uint
	IssueID uint
	Body    string `gorm:"size:2000"`
}

func (r *Reply) SaveReply() error {
	err := database.DB.Create(&r).Error
	return err
}
