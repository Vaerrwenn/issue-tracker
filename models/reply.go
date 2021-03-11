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

// SaveReply saves Reply record to database.
func (r *Reply) SaveReply() error {
	err := database.DB.Create(&r).Error
	return err
}

// FindReplyByID gets a Reply by searching the ID.
func (r *Reply) FindReplyByID(id uint) *Reply {
	var result Reply

	err := database.DB.Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil
	}

	return &result
}

// UpdateReply updates a source Reply.
func (r *Reply) UpdateReply(source *Reply) error {
	err := database.DB.Model(&source).Updates(r).Error
	return err
}
