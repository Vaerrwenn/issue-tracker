package middlewares

import (
	"issue-tracker/auth"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// AuthJWT is a middleware for protected APIs. Checks whether the User who's
// trying to use an API is authenticated or not.
func AuthJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// err := godotenv.Load()
		// if err != nil {
		// 	log.Fatal("Error loading ENV file.")
		// }
		// Get the 'authorization' from the Header.
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "No token in header",
			})
			c.Abort()
			return
		}
		// log.Println(clientToken)
		// // Splits the Bearer and the token. (Used if the token has "Bearer " in front)
		// extractedToken := strings.Split(clientToken, "Bearer ")
		// log.Println(extractedToken)
		// if len(extractedToken) == 2 {
		// 	clientToken = strings.TrimSpace(extractedToken[1])
		// } else {
		// 	c.JSON(http.StatusBadRequest, gin.H{
		// 		"error": "Incorrect format of Authorization token",
		// 	})
		// 	c.Abort()
		// 	return
		// }

		// Check whether the Token is valid.
		jwtWrapper := auth.JwtWrapper{
			SecretKey: os.Getenv("JWT_SECRET"),
			Issuer:    "AuthService",
		}

		claims, err := jwtWrapper.ValidateToken(clientToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Next()
	}
}
