package models

import (
	"issue-tracker/database"

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
