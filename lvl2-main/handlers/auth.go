package handlers

import (
	"fmt"
	"library-management1/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserProfile(c *gin.Context) {
	user, _ := c.Get("currentUser")
	var userData models.User
	userData = user.(models.User)
	fmt.Println(userData.ID)

	c.JSON(http.StatusOK, gin.H{"data": user})
}
