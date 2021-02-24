package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RoleAuth is the middleware to check whether a User with User's Role ID elligible
// to access the request.
//
// Takes allowed Role ID 'allowedRole' (as string) and check it with the Role ID in
// header.
//
// ***ONLY USE THIS WHEN ONLY ONE ROLE IS ALLOWED TO ACCESS THE REQUEST.***
func RoleAuth(allowedRole string) gin.HandlerFunc {
	/*
		I am aware that this function or way is not scalable.
		If this program has more than 2 roles, and several roles can access a request,
		it won't be able to check multiple roles.
		I might have the solution, but will do it if scaling is really required.
	*/
	return func(c *gin.Context) {
		// Get userRole from Header
		headerRole := c.Request.Header.Get("userRole")
		if headerRole == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "No Role ID supplied.",
			})
			c.Abort()
			return
		}

		// Check whether the
		if headerRole != allowedRole {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User with this role is unauthorized to access this Request.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
