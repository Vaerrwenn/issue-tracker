package controllers

import (
	"fmt"
	auth "issue-tracker/auth"
	"issue-tracker/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// RegisterForm binds the data from the Registration Form to the struct.
type RegisterForm struct {
	RoleID   int    `form:"role" json:"role" binding:"required"`
	Name     string `form:"name" json:"name" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// LoginForm binds the data from the Login form to the struct.
type LoginForm struct {
	Email      string `form:"email" json:"email" binding:"required"`
	Password   string `form:"password" json:"password" binding:"required"`
	Remembered bool   `form:"remember" json:"remember"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// RegisterHandler inputs the form-data into the database.
//
// 1. The functions will get the form data.
// If there is an error, the func will send an error to the front.
//
// 2. Password will be encrypted.
//
// 3. Send the data to the model to be saved to the database.
func RegisterHandler(c *gin.Context) {
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

	return
}

// LoginHandler handles the Login feature.
//
// 1. Binds the data from the login form
//
// 2. Check if user with inputted email exists
//
// 3. Check password
//
// 4. Generate JWT Token
//
// 5. Send the Token to the Header.
func LoginHandler(c *gin.Context) {
	// Bind input from the Login form.
	var input LoginForm
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}
	userEmail := models.User{
		Email: input.Email,
	}

	// Check if user with inputted email exists.
	user := userEmail.GetUserByEmail()
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": fmt.Sprintf("Couldn't found user with email: %s", input.Email),
		})
		c.Abort()
		return
	}

	// Check if inputted password is the same as the User's stored password.
	err := bcrypt.CompareHashAndPassword(user.Password, []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user credential",
		})
		c.Abort()
		return
	}

	var expirationHours = 24
	if input.Remembered {
		// Expired in 1 year.
		expirationHours = 8760
	}

	// JWT generation
	jwtWrapper := auth.JwtWrapper{
		SecretKey:       "verysecretkey",
		Issuer:          "AuthService",
		ExpirationHours: int64(expirationHours),
	}

	signedToken, err := jwtWrapper.GenerateToken(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error signing token",
		})
		c.Abort()
		return
	}

	tokenResponse := LoginResponse{
		Token: signedToken,
	}

	c.JSON(http.StatusOK, tokenResponse)

	return
}
