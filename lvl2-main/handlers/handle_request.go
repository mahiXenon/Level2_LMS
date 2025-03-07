package handlers

import (
	"library-management1/database"
	"library-management1/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func HandleRequest(c *gin.Context) {

	user, _ := c.Get("currentUser")
	userData := user.(models.User)
	if userData.Role != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Yoar are not authorized to approve request"})
		return
	}
	var request models.ReaderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var requestEvent models.RequestEvent
	database.DB.Where("id = ?", request.ID).Find(&requestEvent)
	if requestEvent.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request not found"})
		return
	}

	var bookSearch models.BookInventory
	database.DB.Where("isbn = ?", requestEvent.ISBN).Where("library_id = ?", requestEvent.LibraryId).Find(&bookSearch)
	if bookSearch.AvailableCopies <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Book is not available"})
		// var updateRequest models.RequestEvent
		database.DB.Model(&models.RequestEvent{}).Where("id = ?", request.ID).Update("request_type", "reject").Update("approver_id", userData.ID)
		return
	}

	if requestEvent.RequestType == "borrow" {
		bookSearch.AvailableCopies = bookSearch.AvailableCopies - 1
		database.DB.Model(&models.BookInventory{}).Where("isbn = ?", requestEvent.ISBN).Where("library_id = ?", requestEvent.LibraryId).
			Update("available_copies", bookSearch.AvailableCopies)
		issueBook := models.IssueBook{
			ISBN:               requestEvent.ISBN,
			LibraryId:          requestEvent.LibraryId,
			UserID:             requestEvent.UserId,
			IssueApproverId:    userData.ID,
			IssueStatus:        "book issued",
			IssueDate:          time.Now(),
			ExpectedReturnDate: time.Now().AddDate(0, 0, 10),
		}
		database.DB.Create(&issueBook)
		database.DB.Model(&models.RequestEvent{}).Where("id = ?", request.ID).Update("approver_id", userData.ID).Update("approve_date", time.Now())
		c.JSON(http.StatusOK, gin.H{"message": "Request is approved"})
	}

	if requestEvent.RequestType == "return" {
		bookSearch.AvailableCopies = bookSearch.AvailableCopies + 1
		database.DB.Model(&models.BookInventory{}).Where("isbn = ?", requestEvent.ISBN).Where("library_id = ?", requestEvent.LibraryId).
			Update("available_copies", bookSearch.AvailableCopies)
		database.DB.Model(&models.RequestEvent{}).Where("id = ?", request.ID).Update("approver_id", userData.ID).Update("request_type", "Accepted").
			Update("approve_date", time.Now())
		// c.JSON(http.StatusOK, gin.H{"message": "Request is approved"})

		// upadte Issue table
		database.DB.Model(&models.IssueBook{}).Where("isbn = ?", requestEvent.ISBN).Where("library_id = ?", requestEvent.LibraryId).
			Where("user_id = ?", requestEvent.UserId).Update("issue_status", "book returned").Update("return_date", time.Now()).
			Update("return_approver_id", userData.ID)

		c.JSON(http.StatusOK, gin.H{"message": "Book is submitted to liberary"})
	}

}
