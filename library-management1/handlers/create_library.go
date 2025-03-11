package handlers

import (
	// "encoding/json"
	"library-management1/database"
	"library-management1/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateLibrary(c *gin.Context) {
	var authLibrary models.AuthLibrary
	err := c.ShouldBindJSON(&authLibrary)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var libraryFound models.Library
	database.DB.Where("name = ?", authLibrary.Name).Find(&libraryFound)

	if libraryFound.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "library Already Exist"})
		return
	}

	library := models.Library{
		Name: authLibrary.Name,
	}

	database.DB.Create(&library)
	user, _ := c.Get("currentUser")

	//took out library data
	var libraryData models.Library
	database.DB.Where("name = ?", library.Name).Find(&libraryData)

	var userData models.User
	userData = user.(models.User)

	userData.Role = "Owner"
	database.DB.Model(models.User{}).Where("id = ?", userData.ID).Update("Role", userData.Role)
	user_library := models.LibraryUser{
		UserId:    userData.ID,
		LibraryId: libraryData.ID,
	}
	database.DB.Create(&user_library)
	var userUpdate models.User
	database.DB.Where("ID = ?", userData.ID).Find(&userUpdate)

	c.Set("currentUser", userUpdate)
	// var userData models.User
	// userBytes, err := json.Marshal(user)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	// if err := json.Unmarshal(userBytes, &userData); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	// userData = user.(models.User)

	c.JSON(http.StatusOK, gin.H{"data": libraryData})
}
