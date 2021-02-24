package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Initialize DB's connection and migration.
func InitializeDB() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	DB = db
}
