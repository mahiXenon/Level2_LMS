package handlers

import (
	"library-management1/database"
	"library-management1/models"
	"library-management1/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// isValidEmail checks if an email is in a valid format

func CreateUser(c *gin.Context) {
	var authInput models.AuthInput
	if err := c.ShouldBindJSON(&authInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !utils.IsValidEmail(authInput.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid email format"})
		return
	}

	// Validate Contact Number
	if !utils.IsValidContactNumber(authInput.ContactNumber) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid contact number"})
		return
	}

	// Validate Password
	if !utils.IsValidPassword(authInput.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Password must be at least 8 characters long and contain a special character"})
		return
	}
	var validatename string = authInput.Name
	var flag bool = true
	for i := 0; i < len(validatename); i++ {
		if validatename[i] >= '0' && validatename[i] <= '9' {
			flag = false
			break
		}
	}

	if flag == false {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Name must not contain number"})
		return
	} else {

		var userFound models.User
		database.DB.Where("email = ?", authInput.Email).Find(&userFound)

		if userFound.ID != 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "User Already Exist"})
			return
		}

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(authInput.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user := models.User{
			Name:          authInput.Name,
			Email:         authInput.Email,
			Password:      string(passwordHash),
			ContactNumber: authInput.ContactNumber,
			Role:          "user",
		}
		authInput.Password = string(passwordHash)
		database.DB.Create(&user)
		c.JSON(http.StatusOK, gin.H{"data": authInput})
	}

}
