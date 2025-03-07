package handlers

import (
	"fmt"
	"library-management1/database"
	"library-management1/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SearchBook(c *gin.Context) {
	var allBooks []models.BookInventory
	search := c.Param("search")
	fmt.Println(search)
	if search == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	database.DB.Model(&models.BookInventory{}).Where("title LIKE ?", "%"+search+"%").Find(&allBooks)

	c.JSON(http.StatusOK, gin.H{"data": allBooks})
}
