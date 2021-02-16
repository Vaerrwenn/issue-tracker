package models

import "gorm.io/gorm"

type Reply struct {
	gorm.Model
	UserID  int
	User    User
	IssueID int
	Issue   Issue
	Body    string `gorm:"size:2000"`
}
