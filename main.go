package main

import (
	"issue-tracker/controllers"
	"issue-tracker/database"
	"issue-tracker/middlewares"
	migrations "issue-tracker/migrations"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Loads variables on .env file as the machine's local environment if GIN_MODE is not "release".

	mode := os.Getenv("GIN_MODE")

	if mode != "release" {
		err := godotenv.Load()
		if err != nil {
			log.Printf(err.Error())
		}
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
		public := v1.Group("/public")
		{
			public.POST("/register", controllers.RegisterHandler)
			public.POST("/login", controllers.LoginHandler)
		}

		protected := v1.Group("/protected")
		protected.Use(middlewares.AuthJWT())
		{
			issue := protected.Group("/issue")
			{
				issue.POST("/create", middlewares.RoleAuth("1"), controllers.CreateIssueHandler)
				issue.GET("/index", controllers.IndexIssueHandler)
				issue.GET("/show/:id", controllers.ShowIssueHandler)
				issue.PATCH("/update/:id", controllers.UpdateIssueHandler)
				issue.DELETE("/delete/:id", controllers.DeleteIssueHandler)
				issue.POST("/show/:id/reply", controllers.CreateReplyHandler)
				issue.PATCH("/show/:id/update-reply/:replyId", controllers.UpdateReplyHandler)
			}
		}

	}

	r.Run(":" + port)
}
