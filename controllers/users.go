package controllers

import (
	"fmt"
	auth "issue-tracker/auth"
	"issue-tracker/models"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// RegisterForm binds the data from the Registration Form to the struct.
type RegisterForm struct {
	RoleID   int    `form:"role" binding:"required"`
	Name     string `form:"name" binding:"required"`
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

// LoginForm binds the data from the Login form to the struct.
type LoginForm struct {
	Email      string `form:"email" binding:"required"`
	Password   string `form:"password" binding:"required"`
	Remembered bool   `form:"remember"`
}

// ChangePasswordForm is for binding the data from Update Password form.
type ChangePasswordForm struct {
	OldPassword     string `form:"old-password" binding:"required"`
	NewPassword     string `form:"new-password" binding:"required"`
	ConfirmPassword string `form:"confirm-password" binding:"required"`
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
	// Check whether user is logged in.
	token := c.Request.Header.Get("token")
	if token != "" {
		returnErrorAndAbort(c, http.StatusForbidden, "User is already logged in.")
		return
	}
	// Binds the form-data to `input` variable
	var input RegisterForm
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}

	// Password encryption using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		c.Abort()
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
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "User registered successfully.",
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
	// Check whether user is logged in.
	token := c.Request.Header.Get("token")
	if token != "" {
		returnErrorAndAbort(c, http.StatusForbidden, "User is already logged in.")
		return
	}
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
		SecretKey:       os.Getenv("JWT_SECRET"),
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

	c.JSON(http.StatusOK, gin.H{
		"token":     signedToken,
		"userId":    user.ID,
		"userRole":  user.RoleID,
		"userEmail": user.Email,
		"userName":  user.Name,
	})

	return
}

// ChangePasswordHandler handles user's change password request.
func ChangePasswordHandler(c *gin.Context) {
	/*
		Flow:
		1. Validate whether New Password and Confirm Password is the same.
		2. Get User ID from header
		3. Get User data by searching ID
		4. Validate Old Password with User's recorded password
		5. Update Password
	*/
	var input ChangePasswordForm
	if err := c.ShouldBind(&input); err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest, err.Error())
		return
	}

	if input.NewPassword != input.ConfirmPassword {
		returnErrorAndAbort(c, http.StatusBadRequest, "Failed to confirm password.")
		return
	}

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest, err.Error())
		return
	}

	// Get User from User ID
	var user models.User
	source := user.GetUserByID(userID)
	if source == nil {
		returnErrorAndAbort(c, http.StatusNotFound, "User not found.")
		return
	}

	// Check if the Old Password is the same with the new one.
	err = bcrypt.CompareHashAndPassword(source.Password, []byte(input.OldPassword))
	if err != nil {
		returnErrorAndAbort(c, http.StatusForbidden, "Old password is invalid.")
		return
	}

	// Generate New Password.
	newPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest, "Failed to encrypt password.")
		return
	}

	err = source.UpdatePassword(newPassword)
	if err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "Password updated successfully.",
	})
	return
}

// ShowUserHandler handles Show User data request.
//
// Requires:
//
// - UserID from the param (URL)
func ShowUserHandler(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest, "No user ID provided.")
		return
	}

	var user models.User
	result := user.GetUserByID(userID)
	if result == nil {
		returnErrorAndAbort(c, http.StatusNotFound, "No User found.")
		return
	}
	result.Password = nil

	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}
