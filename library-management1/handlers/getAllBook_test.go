package handlers

import (
	"encoding/json"
	"library-management1/database"
	"library-management1/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Setup router for testing
func SetupGetAllBookRouter() *gin.Engine {
	r := gin.Default()
	r.Use(MockMiddlewareForBookRetrieval()) // Middleware to set user context
	r.GET("/books", GetAllBook)
	return r
}

// MockMiddlewareForBookRetrieval to simulate authentication and set user context
func MockMiddlewareForBookRetrieval() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		database.DB.Where("email = ?", "user@example.com").First(&user)
		if user.ID != 0 {
			c.Set("currentUser", user)
		}
		c.Next()
	}
}

func TestGetAllBook(t *testing.T) {
	database.DB = database.SetDB()
	r := SetupGetAllBookRouter()

	// Clear database before testing
	database.DB.Exec("DELETE FROM book_inventories")
	database.DB.Exec("DELETE FROM library_users")
	database.DB.Exec("DELETE FROM libraries")
	database.DB.Exec("DELETE FROM users")

	// Create a user
	user := models.User{
		Name:          "Test User",
		Email:         "user@example.com",
		Password:      "Password@123",
		ContactNumber: "1234567890",
		Role:          "user",
	}
	database.DB.Create(&user)

	// Test case: User not registered in any library
	t.Run("User not registered in library", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/books", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var responseBody map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.Equal(t, "You have not registered your library yet", responseBody["message"])
	})

	// Create a library and register the user
	library := models.Library{
		Name: "Central Library",
	}
	database.DB.Create(&library)

	libraryUser := models.LibraryUser{
		UserId:    user.ID,
		LibraryId: library.ID,
	}
	database.DB.Create(&libraryUser)

	// Test case: No books available in library
	t.Run("No books in library", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/books", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var responseBody map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.Equal(t, []interface{}{}, responseBody["data"]) // Expect empty list
	})

	// Insert book into library
	book := models.BookInventory{
		ISBN:            "978-3-16-148410-0",
		LibraryId:       library.ID,
		Title:           "Sample Book",
		Author:          "Author Name",
		Publisher:       "Publisher Name",
		Version:         "1st Edition",
		TotalCopies:     5,
		AvailableCopies: 5,
	}
	database.DB.Create(&book)

	// Test case: Successfully retrieve books
	t.Run("Retrieve books from library", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/books", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// var responseBody map[string]interface{}
		// json.Unmarshal(w.Body.Bytes(), &responseBody)

		// // Check that books are returned
		// books, exists := responseBody["data"].([]interface{})
		// assert.True(t, exists, "Data should be a list of books")
		// assert.Len(t, books, 1, "There should be 1 book in response")
	})
}
