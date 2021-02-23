package middlewares

import (
	"issue-tracker/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthJWT is a middleware for protected APIs. Checks whether the User who's
// trying to use an API is authenticated or not.
func AuthJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the 'authorization' from the Header.
		clientToken := c.Request.Header.Get("authorization")
		if clientToken == "" {
			c.JSON(http.StatusForbidden, "No authorization header provided")
			c.Abort()
			return
		}

		// Gets the token
		extractedToken := strings.Split(clientToken, "Bearer ")
		if len(extractedToken) == 2 {
			clientToken = strings.TrimSpace(extractedToken[1])
		} else {
			c.JSON(http.StatusBadRequest, "Incorrect format of Authorization token")
			c.Abort()
			return
		}

		// Check whether the Token is valid.
		jwtWrapper := auth.JwtWrapper{
			SecretKey: "verysecretkey",
			Issuer:    "AuthService",
		}

		claims, err := jwtWrapper.ValidateToken(clientToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Next()
	}
}
