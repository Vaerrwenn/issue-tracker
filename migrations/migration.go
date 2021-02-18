package migrations

import (
	"issue-tracker/models"

	"gorm.io/gorm"
)

// MigrateTables migrates the Models into the Database table.
func MigrateTables(db *gorm.DB) {
	db.AutoMigrate(&models.Role{},
		&models.User{},
		&models.Issue{},
		&models.Reply{},
		&models.Notification{},
	)
}
