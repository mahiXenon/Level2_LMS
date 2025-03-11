package handlers

import (
	"library-management1/database"
	"library-management1/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllBook(c *gin.Context) {
	var allBooks []models.BookInventory

	if err := database.DB.Find(&allBooks).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, _ := c.Get("currentUser")
	userData := user.(models.User)
	var library models.LibraryUser
	database.DB.Where("user_id = ?", userData.ID).Find(&library)

	if library.LibraryId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "You have not registered your library yet"})
		return
	}

	database.DB.Model(&models.BookInventory{}).Where("library_id = ?", library.LibraryId).Find(&allBooks)

	c.JSON(http.StatusOK, gin.H{"data": allBooks})
}
