package handlers

import (
	"library-management1/database"
	"library-management1/models"

	"github.com/gin-gonic/gin"
)

var users models.User

func Test(c *gin.Context) {
	// user := models.User{Name: "Nov Nov", Email: "nov@example.com"}
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(500, err.Error())
		return
	}
	database.DB.Create(&user)
	users = user
}

func Del(c *gin.Context) {
	database.DB.Delete(&users)

}
