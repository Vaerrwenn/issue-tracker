package controllers

import (
	"issue-tracker/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// RegisterForm binds the data from the Form to the struct.
type RegisterForm struct {
	RoleID   int    `form:"role" json:"role" binding:"required"`
	Name     string `form:"name" json:"name" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type LoginForm struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json"password" binding:"required"`
}

// CreateUser inputs the form-data into the database.
// Password will be encrypted here.
func CreateUser(c *gin.Context) {
	// Binds the form-data to `input` variable
	var input RegisterForm
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Password encryption using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Assign the input into user variable. Used for storing the data in database.
	user := models.User{
		RoleID:   input.RoleID,
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
	}
	// Saves user data.
	err = user.SaveUserData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": 1,
	})
}
