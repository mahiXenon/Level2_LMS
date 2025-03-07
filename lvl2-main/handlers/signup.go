package handlers

import (
	"library-management1/database"
	"library-management1/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *gin.Context) {
	var authInput models.AuthInput
	if err := c.ShouldBindJSON(&authInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

		database.DB.Where("contact_number = ?", authInput.ContactNumber).Find(&userFound)
		if userFound.ContactNumber != "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Contact Number Already Exist"})
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
