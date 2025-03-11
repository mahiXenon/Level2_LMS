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

func SetupSearchBookRouter() *gin.Engine {
	r := gin.Default()
	r.Use(MockMiddlewareForSearchRouter())
	r.GET("/search-book/:search", SearchBook)
	return r
}
func MockMiddlewareForSearchRouter() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		database.DB.Where("email = ?", "user@example.com").First(&user)
		if user.ID != 0 {
			c.Set("currentUser", user)
		}
		c.Next()
	}
}

func TestSearchBook(t *testing.T) {
	r := SetupSearchBookRouter()
	database.DB = database.SetDB()

	// Clear existing books
	database.DB.Exec("DELETE FROM book_inventories")
	database.DB.Exec("DELETE FROM users")

	user := models.User{

		Name:          "Test User",
		Email:         "user@example.com",
		Password:      "Password@123",
		ContactNumber: "1234567890",
		Role:          "admin",
	}
	database.DB.Create(&user)

	// Insert test books
	books := models.BookInventory{
		ISBN: "12345", Title: "Go Programming", Author: "John Doe", Publisher: "TechBooks", Version: "1st", TotalCopies: 5, AvailableCopies: 5,
		// {ISBN: "67890", Title: "Python Basics", Author: "Jane Smith", Publisher: "CodeWorld", Version: "2nd", TotalCopies: 3, AvailableCopies: 3},
	}
	books1 := models.BookInventory{
		ISBN: "54321", Title: "The Sql Syntax", Author: "Jatin", Publisher: "TechBooks", Version: "1st", TotalCopies: 5, AvailableCopies: 5,
	}
	database.DB.Create(&books)
	database.DB.Create(&books1)

	testCases := []struct {
		name         string
		searchQuery  string
		expectedCode int
		expectedLen  int
	}{
		{
			name:         "Search with matching result",
			searchQuery:  "Go Programming",
			expectedCode: http.StatusOK,
			expectedLen:  1,
		},
		{
			name:         "Search with no matching result",
			searchQuery:  "Rust Programming",
			expectedCode: http.StatusOK,
			expectedLen:  0,
		},
		// {
		// 	name:         "Empty search query",
		// 	searchQuery:  "",
		// 	expectedCode: http.StatusBadRequest,
		// 	expectedLen:  -1, // This case won't check length, just error message
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/search-book/"+tc.searchQuery, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)

			if tc.expectedLen >= 0 {
				var response map[string][]models.BookInventory
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Len(t, response["data"], tc.expectedLen)
			}
		})
	}
}
