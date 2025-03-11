package handlers

import (
	"library-management1/database"
	"library-management1/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllLibrary(c *gin.Context) {
	var alllibrary []models.Library

	database.DB.Find(&alllibrary)

	c.JSON(http.StatusOK, gin.H{"data": alllibrary})
}
