package handlers

import (
	"library-management1/database"
	"library-management1/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func InsertBook(c *gin.Context) {
	var inputBook models.InputBook
	if err := c.ShouldBindJSON(&inputBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, _ := c.Get("currentUser")
	userData := user.(models.User)
	if userData.Role != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Only admin can insert book"})
		return
	}
	if userData.Role == "admin" {
		var library models.LibraryUser
		// finding the library id of user admin
		database.DB.Where("user_id = ?", userData.ID).Find(&library)
		// finding isbn and library id in book inventory
		var checkBook models.BookInventory
		database.DB.Where("isbn = ?", inputBook.ISBN).Where("library_id = ?", library.LibraryId).Find(&checkBook)
		if checkBook.ISBN != "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Book is already present in library"})
			return
		}
		book := models.BookInventory{
			ISBN:            inputBook.ISBN,
			LibraryId:       library.LibraryId,
			Title:           inputBook.Title,
			Author:          inputBook.Author,
			Publisher:       inputBook.Publisher,
			Version:         inputBook.Version,
			TotalCopies:     inputBook.TotalCopies,
			AvailableCopies: inputBook.TotalCopies,
		}
		inputBook.Title = strings.ToUpper(inputBook.Title)
		inputBook.Author = strings.ToUpper(inputBook.Author)
		inputBook.Publisher = strings.ToUpper(inputBook.Publisher)
		if book.TotalCopies < 0 {
			c.JSON(http.StatusBadGateway, gin.H{"message": "available_copies or total_copies can not be negative"})
			return
		}
		database.DB.Create(&book)
		c.JSON(http.StatusOK, gin.H{"data": book})
	}

}
