package main

import (
	controllers "issue-tracker/controllers"
	"issue-tracker/database"
	migrations "issue-tracker/migrations"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Loads variables on .env file as the machine's local environment.
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Gets PORT on environment.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize Database.
	database.InitializeDB()

	// Table Migration
	migrations.MigrateTables(database.DB)

	// Initiate Gin's default engine
	r := gin.Default()

	v1 := r.Group("/v1")
	{
		v1.POST("/register", controllers.CreateUser)
		v1.POST("/login")
	}

	r.Run(":" + port)
}
