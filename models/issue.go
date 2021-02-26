package models

import (
	"errors"
	"fmt"
	"issue-tracker/database"
	"strconv"
	"time"

	"gorm.io/gorm"
)

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

// IssueIndex is used for IndexIssue operation.
type IssueIndex struct {
	ID        int
	Title     string
	Status    string
	Severity  string
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    int
	UserName  string
}

// IssueShow is used for ShowIssue operation.
type IssueShow struct {
	ID        int
	Title     string
	Body      string
	Status    string
	Severity  string
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    int
	UserName  string
}

// ValidateIssue validates the Issue data.
// Some validation is done by GORM, but there are other validation,
// id est Severity being only 1 - 3, needs to be validated manually.
func (i *Issue) ValidateIssue() error {
	severity, err := strconv.Atoi(i.Severity)
	if err != nil {
		return err
	}
	if severity > 3 || severity < 1 {
		return errors.New("ERROR SEVERITY: severity must be between 1 - 3")
	}
	return nil
}

// SaveIssue saves the Issue to the database.
func (i *Issue) SaveIssue() error {
	err := database.DB.Create(&i).Error
	return err
}

// IndexIssues fetches all issues from the database.
func (i *Issue) IndexIssues() (*[]IssueIndex, error) {
	var issues []IssueIndex
	query := database.DB.Model(&Issue{}).
		Select(`
			issues.id,
			issues.title,
			issues.status,
			issues.severity,
			issues.created_at,
			issues.updated_at,
			issues.user_id,
			users."name" AS "user_name"`).
		Joins("left join users on issues.user_id = users.id").
		Scan(&issues)

	if query.Error != nil {
		return nil, query.Error
	}

	return &issues, nil
}

// FindIssueByID fetches an issue with provided ID.
func (i *Issue) FindIssueByID(id string) (*IssueShow, error) {
	var result IssueShow
	query := database.DB.Model(&Issue{}).
		Select(`
			issues.id,
			issues.title,
			issues.body,
			issues.status,
			issues.severity,
			issues.created_at,
			issues.updated_at,
			issues.user_id,
			users."name" AS "user_name"`).
		Joins("left join users on issues.user_id = users.id").
		Where("issues.id = ?", id).
		Scan(&result)

	if result.ID == 0 {
		return nil, fmt.Errorf("ERROR: could not find issue with ID: %s", id)
	}

	if query.Error != nil {
		return nil, query.Error
	}
	return &result, nil
}
