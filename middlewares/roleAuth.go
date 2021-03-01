package middlewares

import (
	"issue-tracker/models"
	"net/http"
	"strconv"

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
		var user models.User
		userID, err := strconv.Atoi(c.Request.Header.Get("userID"))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		userRole, err := user.GetUserRoleByID(userID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if userRole != allowedRole {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User is unauthorized to use this request."})
			c.Abort()
			return
		}

		c.Next()
	}
}
