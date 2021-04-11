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
	UserID            int
	Title             string `gorm:"size:100"`
	Body              string `gorm:"size:2000"`
	Status            string `gorm:"size:1"` // 1 = Opened or 0 = Closed
	Severity          string `gorm:"size:1"` // 1 = Low, 2 = Medium, 3 = High
	UpdatedByUserID   int
	UpdatedByUserName string `gorm:"size:100"`
	Replies           []Reply
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

type IssueShow struct {
	ID        int
	Title     string
	Status    string
	Severity  string
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    int
	UserName  string
}

type RepliesInIssue struct {
	ID        int
	UserID    int
	IssueID   int
	Body      string
	Replier   string
	CreatedAt time.Time
	UpdatedAt time.Time
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

// FindIssueAndRepliesByID fetches an issue with provided ID.
// It will return issue, replies of that issue, and user data that is needed for Show route.
func (i *Issue) FindIssueAndRepliesByID(id string) (*IssueShow, *[]RepliesInIssue, error) {
	var issue IssueShow
	// query := database.DB.Preload("Replies").Where("issues.id = ?", id).First(&result)
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
		Where("issues.id = ?", id).
		First(&issue)

	var replies []RepliesInIssue
	queryReplies := database.DB.Model(&Reply{}).
		Select(`
			replies.id,
			replies.user_id,
			replies.issue_id,
			users."name" as "replier",
			replies.body,
			replies.created_at,
			replies.updated_at`).
		Joins("join users on replies.user_id = users.id").
		Joins("join issues on replies.issue_id = issues.id").
		Where("replies.issue_id = ?", id).
		Scan(&replies)

	if issue.ID == 0 {
		return nil, nil, fmt.Errorf("ERROR: could not find issue with ID: %s", id)
	}

	if query.Error != nil || queryReplies.Error != nil {
		return nil, nil, query.Error
	}
	return &issue, &replies, nil
}

func (i *Issue) FindOneIssueByID(id string) (*Issue, error) {
	var result Issue
	query := database.DB.Preload("Replies").Where("issues.id = ?", id).First(&result)

	if result.ID == 0 {
		return nil, fmt.Errorf("ERROR: Could not find issue with ID: %s", id)
	}

	if query.Error != nil {
		return nil, query.Error
	}
	return &result, nil
}

// UpdateIssue updates an Issue data.
// Takes an origin Issue as parameter. Origin issue is
// the Issue that user will update.
func (i *Issue) UpdateIssue(origin *Issue) error {
	err := database.DB.Model(&origin).Updates(i).Error
	return err
}

// DeleteIssue deletes an Issue data.
func (i *Issue) DeleteIssue() error {
	err := database.DB.Delete(&i).Error
	return err
}
