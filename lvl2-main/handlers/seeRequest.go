package handlers

import (
	"fmt"
	"library-management1/database"
	"library-management1/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SeeRequest(c *gin.Context) {

	user, _ := c.Get("currentUser")
	userData := user.(models.User)
	if userData.Role != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Only admin can see all request"})
		return
	}
	if userData.Role == "admin" {
		var library models.LibraryUser
		// finding the library id of user admin
		database.DB.Where("user_id = ?", userData.ID).Find(&library)
		fmt.Println("lib", library.LibraryId)
		// finding all request of library
		var listRequest []models.RequestEvent
		database.DB.Model(&models.RequestEvent{}).Where("library_id = ?", library.LibraryId).Find(&listRequest)

		if len(listRequest) != 0 {

			var request []models.ListRequest
			for i := 0; i < len(listRequest); i++ {
				list := models.ListRequest{
					ID:          listRequest[i].ID,
					ISBN:        listRequest[i].ISBN,
					UserId:      listRequest[i].UserId,
					RequestDate: listRequest[i].RequestDate,
					RequestType: listRequest[i].RequestType,
				}
				request = append(request, list)

			}
			fmt.Println("request", request)
			c.JSON(http.StatusOK, gin.H{"data": request})
		}
	}
}
