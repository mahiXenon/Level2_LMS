package handlers

import (
	"library-management1/database"
	"library-management1/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestEvent(c *gin.Context) {
	var inputRequest models.RequestInput
	if err := c.ShouldBindJSON(&inputRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, _ := c.Get("currentUser")
	userData := user.(models.User)

	var libraryPresent models.Library
	database.DB.Where("id = ?", inputRequest.LibraryId).Find(&libraryPresent)
	if libraryPresent.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Library is not present"})
		return
	}

	var library models.LibraryUser
	database.DB.Where("user_id = ?", userData.ID).Where("library_id = ?", inputRequest.LibraryId).Find(&library)
	if library.UserId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User is not authorized to request book user must be register in library"})
		return
	}
	var checkBook models.BookInventory
	database.DB.Where("isbn = ?", inputRequest.ISBN).Where("library_id = ?", inputRequest.LibraryId).Find(&checkBook)
	if checkBook.ISBN == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Book is not present in library"})
		return
	}

	Request := models.RequestEvent{
		ISBN:        inputRequest.ISBN,
		UserId:      userData.ID,
		LibraryId:   inputRequest.LibraryId,
		RequestDate: time.Now(),
		RequestType: inputRequest.RequestType,
	}
	var alreadyRequested models.RequestEvent
	database.DB.Where("isbn = ?", inputRequest.ISBN).Where("user_id = ?", userData.ID).
		Where("library_id = ?", inputRequest.LibraryId).Where("request_type = ?", Request.RequestType).Find(&alreadyRequested)
	if alreadyRequested.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Book is already requested"})
		return
	}
	if inputRequest.RequestType == "borrow" {
		database.DB.Create(&Request)
		c.JSON(http.StatusOK, gin.H{"message": "Request is successfully placed"})
		return
	} else {
		var requestSearch models.RequestEvent
		database.DB.Where("isbn = ?", inputRequest.ISBN).Where("user_id = ?", userData.ID).
			Where("library_id = ?", inputRequest.LibraryId).Where("request_type = ?", "borrow").Find(&requestSearch)

		if requestSearch.ID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "You have not taken any book"})
			return
		}
		database.DB.Model(&models.RequestEvent{}).Where("isbn = ?", inputRequest.ISBN).Where("user_id = ?", userData.ID).
			Where("library_id = ?", inputRequest.LibraryId).Where("request_type = ?", "borrow").Update("request_type", "return")

		c.JSON(http.StatusOK, gin.H{"message": "Returned Book Successfully"})
		return

	}

}
