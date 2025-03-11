package handlers

import (
	"bytes"
	"encoding/json"
	"library-management1/database"
	"library-management1/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetupInsertBookRouter() *gin.Engine {
	r := gin.Default()
	r.Use(MockMiddlewareForBookInsertion()) // Middleware to set user context
	r.POST("/insert-book", InsertBook)
	return r
}

// MockMiddlewareForBookInsertion to simulate authentication and set user context
func MockMiddlewareForBookInsertion() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		database.DB.Where("email = ?", "admin@example.com").First(&user)
		if user.ID != 0 {
			c.Set("currentUser", user)
		}
		if user.ID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func TestInsertBook(t *testing.T) {
	database.DB = database.SetDB()
	r := SetupInsertBookRouter()

	// Clear database before testing
	database.DB.Exec("DELETE FROM book_inventories")
	database.DB.Exec("DELETE FROM library_users")
	database.DB.Exec("DELETE FROM libraries")
	database.DB.Exec("DELETE FROM users")

	// Create an admin user
	admin := models.User{
		Name:          "Admin User",
		Email:         "admin@example.com",
		Password:      "Password@123",
		ContactNumber: "1234567890",
		Role:          "admin",
	}
	database.DB.Create(&admin)

	// Create a library and link the admin to it
	library := models.Library{
		Name: "Main Library",
	}
	database.DB.Create(&library)

	libraryUser := models.LibraryUser{
		UserId:    admin.ID,
		LibraryId: library.ID,
	}
	database.DB.Create(&libraryUser)
	// insert book in book inventory
	bookInsert := models.BookInventory{
		ISBN:            "978-3-16-148410-0",
		LibraryId:       library.ID,
		Title:           "TEST BOOK",
		Author:          "TEST AUTHOR",
		Publisher:       "TEST PUBLISHER",
		Version:         "1st Edition",
		TotalCopies:     10,
		AvailableCopies: 10,
	}
	database.DB.Create(&bookInsert)

	// Define test cases
	testCases := []struct {
		name         string
		input        models.InputBook
		expectedCode int
		expectedBody map[string]interface{}
	}{
		// {
		// 	name: "Successful book insertion",
		// 	input: models.InputBook{
		// 		ISBN:        "978-3-16-148410-0",
		// 		Title:       "TEST BOOK",
		// 		Author:      "TEST AUTHOR",
		// 		Publisher:   "TEST PUBLISHER",
		// 		Version:     "1st Edition",
		// 		TotalCopies: 10,
		// 	},
		// 	expectedCode: http.StatusOK,
		// 	expectedBody: map[string]interface{}{
		// 		"data": map[string]interface{}{
		// 			"ISBN":            "978-3-16-148410-0",
		// 			"LibraryId":       float64(library.ID), // Checking LibraryId
		// 			"Title":           "TEST BOOK",
		// 			"Author":          "TEST AUTHOR",
		// 			"Publisher":       "TEST PUBLISHER",
		// 			"Version":         "1st Edition",
		// 			"TotalCopies":     10,
		// 			"AvailableCopies": 10, // Available copies must match total copies initially
		// 		},
		// 	},
		// },
		{
			name: "Book already exists in library",
			input: models.InputBook{
				ISBN:        "978-3-16-148410-0",
				Title:       "Test Book",
				Author:      "Test Author",
				Publisher:   "Test Publisher",
				Version:     "1st Edition",
				TotalCopies: 5,
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"message": "Book is already present in library",
			},
		},
		{
			name: "Negative book copies",
			input: models.InputBook{
				ISBN:        "978-3-16-148411-7",
				Title:       "New Book",
				Author:      "New Author",
				Publisher:   "New Publisher",
				Version:     "2nd Edition",
				TotalCopies: -5,
			},
			expectedCode: http.StatusBadGateway,
			expectedBody: map[string]interface{}{
				"message": "available_copies or total_copies can not be negative",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonValue, _ := json.Marshal(tc.input)
			req, _ := http.NewRequest("POST", "/insert-book", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)

			// Parse JSON response
			var responseBody map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(t, err, "Response should be valid JSON")

			// Validate response message
			for key, expectedValue := range tc.expectedBody {
				assert.Contains(t, responseBody, key, "Response should contain key: "+key)

				// Check nested data structure if exists
				if key == "data" {
					dataMap, ok := responseBody["data"].(map[string]interface{})
					assert.True(t, ok, "Response data should be a map")
					for dataKey, dataExpectedValue := range expectedValue.(map[string]interface{}) {
						assert.Equal(t, dataExpectedValue, dataMap[dataKey], "Mismatch in response for key: "+dataKey)
					}
				} else {
					assert.Equal(t, expectedValue, responseBody[key], "Mismatch in response for key: "+key)
				}
			}
		})
	}
}
