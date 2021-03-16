package models

import (
	"issue-tracker/database"
	"strconv"

	"gorm.io/gorm"
)

// User belongs to a Role
type User struct {
	gorm.Model
	RoleID        int
	Role          Role
	Issues        []Issue
	Replies       []Reply
	Notifications []Notification
	Name          string `gorm:"size:100"`
	Email         string `gorm:"size:300"`
	Password      []byte
}

// SaveUserData saves a User's data from Register.
// Returns error if failed.
func (u *User) SaveUserData() error {
	err := database.DB.Create(&u).Error
	return err
}

// GetUserByEmail searches a User by presented Email.
// Returns the User data.
func (u *User) GetUserByEmail() *User {
	var result = &User{}
	err := database.DB.Where(map[string]interface{}{
		"email": u.Email,
	}).First(&result).Error
	if err != nil {
		return nil
	}
	return result
}

// GetUserRoleByID fetches a User's RoleID by searching its ID.
// Returns a string and error.
func (u *User) GetUserRoleByID(id int) (string, error) {
	var result string
	var user User

	query := database.DB.Where("id = ?", id).First(&user)
	if query.Error != nil {
		return "", query.Error
	}

	result = strconv.Itoa(user.RoleID)
	return result, nil
}

// GetUserByID gets a User data by ID.
func (u *User) GetUserByID(id int) *User {
	var result User
	err := database.DB.Joins("Role").Preload("Issues").Preload("Replies").Where("users.id = ?", id).First(&result).Error
	if err != nil {
		return nil
	}
	return &result
}

// UpdatePassword updates a User's password.
func (u *User) UpdatePassword(newPassword []byte) error {
	err := database.DB.Model(&u).Update("password", newPassword).Error
	return err
}
