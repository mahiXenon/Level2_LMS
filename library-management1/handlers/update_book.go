package handlers

import (
	"library-management1/database"
	"library-management1/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpadateBookCopies(c *gin.Context) {
	var updateBook models.UpdateBookDetails
	if err := c.ShouldBindJSON(&updateBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, _ := c.Get("currentUser")
	userData := user.(models.User)
	if userData.Role != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Only admin can update book"})
		return
	}

	if userData.Role == "admin" {
		var library models.LibraryUser
		// finding the library id of user admin
		database.DB.Where("user_id = ?", userData.ID).Find(&library)
		// finding isbn and library id in book inventory
		var checkBook models.BookInventory
		database.DB.Where("isbn = ?", updateBook.ISBN).Where("library_id = ?", library.LibraryId).Find(&checkBook)
		if checkBook.ISBN == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Book is not present in library"})
			return
		}
		checkBook.TotalCopies = updateBook.ADD + checkBook.TotalCopies
		checkBook.TotalCopies = checkBook.TotalCopies - updateBook.DecreaseCount
		checkBook.AvailableCopies = checkBook.AvailableCopies + updateBook.ADD
		checkBook.AvailableCopies = checkBook.AvailableCopies - updateBook.DecreaseCount
		if checkBook.TotalCopies < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Decrease count is more than total copies"})
			return
		}
		if checkBook.AvailableCopies < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Decrease count is more than available copies"})
			return
		}
		if checkBook.TotalCopies == 0 {
			// delete book from inventory
			database.DB.Where("isbn = ?", updateBook.ISBN).Where("library_id = ?", library.LibraryId).Delete(&models.BookInventory{})
			return
		}
		database.DB.Model(models.BookInventory{}).Where("isbn = ?", updateBook.ISBN).Where("library_id = ?", library.LibraryId).
			Update("total_copies", checkBook.TotalCopies).Update("available_copies", checkBook.AvailableCopies)
		var dislay models.BookInventory
		database.DB.Where("isbn = ?", updateBook.ISBN).Where("library_id = ?", library.LibraryId).Find(&dislay)
		c.JSON(http.StatusOK, gin.H{"data": dislay})

	}
}
