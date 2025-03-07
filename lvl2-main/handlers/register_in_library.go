package handlers

import (
	"library-management1/database"
	"library-management1/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	user, _ := c.Get("currentUser")
	var userData models.User
	userData = user.(models.User)

	var register models.RegisterLibrary
	if err := c.ShouldBindJSON(&register); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if userData.Role == "admin" {
		var library models.LibraryUser
		// finding the library id of user admin
		database.DB.Where("user_id = ?", userData.ID).Find(&library)

		var searchLibrary models.Library
		database.DB.Where("id = ?", library.LibraryId).Find(&searchLibrary)
		if searchLibrary.Name != register.Name {
			c.JSON(http.StatusBadRequest, gin.H{"message": "You are not allowed to register in other library"})
			return
		}
	}

	var library models.Library
	database.DB.Where("name = ?", register.Name).Find(&library)
	if library.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No library found"})
		return
	}

	library_users := models.LibraryUser{
		UserId:    userData.ID,
		LibraryId: library.ID,
	}

	database.DB.Create(&library_users)
	c.JSON(http.StatusOK, gin.H{"message": "Sucessfully Registered"})

	// abhi or karna hai

}
