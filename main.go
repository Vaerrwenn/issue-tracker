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

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"data": "There is nothing here."})
	})

	r.GET("/robots", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "Please don't do anything bad to this service :)."})
	})

	v1 := r.Group("/v1")
	{
		// All requests on "/public" does not require any Header.
		public := v1.Group("/public")
		{
			public.POST("/register", controllers.RegisterHandler)
			public.POST("/login", controllers.LoginHandler)
		}

		// All requests in protected requires at least:
		// - token in Header
		// - userID in Header
		protected := v1.Group("/protected")
		protected.Use(middlewares.AuthJWT())
		{
			// NOTE:
			// The :id Param from here and beyond are the ID of each Router group's (user or issue).

			user := protected.Group("/user")
			{
				// Only requires the Param :id from URL.
				user.PATCH("/:id/change-password", controllers.ChangePasswordHandler)
				// Only requires the Param :id from URL.
				user.GET("/:id", controllers.ShowUserHandler)
			}

			issue := protected.Group("/issue")
			{
				// Require form data with input name as:
				// - tite
				// - description
				// - severity
				issue.POST("/create", middlewares.RoleAuth("1"), controllers.CreateIssueHandler)

				// No other requirement needed.
				issue.GET("/index", controllers.IndexIssueHandler)

				// Only requires the Param :id from URL.
				issue.GET("/show/:id", controllers.ShowIssueHandler)

				// Requires:
				// - Param :id from URL
				// - Form with input name as follows:
				//     - title
				//     - description
				//     - status
				//     - severity
				// - userID from Header
				issue.PATCH("/update/:id", controllers.UpdateIssueHandler)

				// Requires:
				// - Param :id from URL
				// - userID from Header
				issue.DELETE("/delete/:id", middlewares.RoleAuth("1"), controllers.DeleteIssueHandler)

				// Requires:
				// - Param :id from URL
				// - userID from Header
				// - Form with input name as follows:
				// 	   - description
				issue.POST("/show/:id/reply", controllers.CreateReplyHandler)

				// Requires:
				// - Param :id from URL
				// - Param :replyId from URL
				// - userID from Header
				// - Form with input name as follows:
				// 	   - description
				issue.PATCH("/show/:id/update-reply/:replyId", controllers.UpdateReplyHandler)

				// Requires:
				// - Param :id from URL
				// - Param :replyId from URL
				// - userID from Header
				issue.DELETE("/show/:id/delete-reply/:replyId", controllers.DeleteReplyHandler)
			}
		}

	}

	r.Run(":" + port)
}
