package main

import (
	"issue-tracker/controllers"
	"issue-tracker/database"
	"issue-tracker/middlewares"
	migrations "issue-tracker/migrations"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
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

	// CORS stuff.
	corsConfig := cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
	r.Use(corsConfig)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, "We're no strangers to love || You know the rules and so do I || "+
			"A full commitment's what I'm thinking of || You wouldn't get this from any other guy || "+
			"I just wanna tell you how I'm feeling || Gotta make you understand || "+
			"Never gonna give you up || Never gonna let you down || Never gonna run around and desert you || "+
			"Never gonna make you cry || Never gonna say goodbye || Never gonna tell a lie and hurt you || ",
		)
		return
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
