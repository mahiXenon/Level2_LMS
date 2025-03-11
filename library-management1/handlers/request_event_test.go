package handlers

// import (
// 	"bytes"
// 	"encoding/json"
// 	"library-management1/database"
// 	"library-management1/models"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/assert"
// )

// // Setup router for testing
// func SetupRequestEventRouter() *gin.Engine {
// 	r := gin.Default()
// 	r.Use(MockMiddlewareForEventRouter()) // Middleware to set user context
// 	r.POST("/request-event", RequestEvent)
// 	return r
// }

// // MockMiddlewareForEventRouter simulates authentication and sets user context
// func MockMiddlewareForEventRouter() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var user models.User
// 		database.DB.Where("email = ?", "user@example.com").First(&user)
// 		if user.ID != 0 {
// 			c.Set("currentUser", user)
// 		}
// 		c.Next()
// 	}
// }

// func TestRequestEventCases(t *testing.T) {
// 	r := SetupRequestEventRouter()
// 	database.DB = database.SetDB()

// 	database.DB.Exec("DELETE FROM libraries")
// 	database.DB.Exec("DELETE FROM users")
// 	database.DB.Exec("DELETE FROM library_users")
// 	database.DB.Exec("DELETE FROM book_inventories")
// 	database.DB.Exec("DELETE FROM request_events")
// 	library := models.Library{
// 		ID:   20,
// 		Name: "Test Library",
// 	}
// 	database.DB.Create(&library) // Insert the library

// 	newlibrary := models.Library{
// 		ID:   21,
// 		Name: "Test Library",
// 	}
// 	database.DB.Create(&newlibrary) // Insert the library

// 	// inser a book

// 	book := models.BookInventory{
// 		ISBN:            "978-3-16-148410-0",
// 		LibraryId:       20,
// 		Title:           "Book Title",
// 		Author:          "Author",
// 		Publisher:       "Publisher",
// 		Version:         "1.0",
// 		TotalCopies:     10,
// 		AvailableCopies: 10,
// 	}
// 	database.DB.Create(&book) // Insert the book

// 	secondbook := models.BookInventory{
// 		ISBN:            "978-3-16-148411-1",
// 		LibraryId:       20,
// 		Title:           "Book Title",
// 		Author:          "Author",
// 		Publisher:       "Publisher",
// 		Version:         "1.0",
// 		TotalCopies:     10,
// 		AvailableCopies: 10,
// 	}
// 	database.DB.Create(&secondbook)

// 	// insert third book
// 	thirdbook := models.BookInventory{
// 		ISBN:            "978-3-16-148412-1",
// 		LibraryId:       20,
// 		Title:           "Book Title",
// 		Author:          "Author",
// 		Publisher:       "Publisher",
// 		Version:         "1.0",
// 		TotalCopies:     10,
// 		AvailableCopies: 10,
// 	}
// 	database.DB.Create(&thirdbook)
// 	// Insert a test user (optional)

// 	user := models.User{

// 		Name:          "Test User",
// 		Email:         "user@example.com",
// 		Password:      "Password@123",
// 		ContactNumber: "1234567890",
// 		Role:          "user",
// 	}
// 	database.DB.Create(&user) // Insert the user

// 	libraryUser := models.LibraryUser{
// 		UserId:    user.ID,
// 		LibraryId: library.ID,
// 	}
// 	database.DB.Create(&libraryUser) // Insert the library user
// 	// insert a request in table for duplicate request
// 	request := models.RequestEvent{
// 		UserId:      user.ID,
// 		ISBN:        "978-3-16-148411-1",
// 		LibraryId:   20,
// 		RequestType: "borrow",
// 	}
// 	database.DB.Create(&request) // Insert the request
// 	// Test cases slice
// 	testCases := []struct {
// 		name         string
// 		requestBody  models.RequestInput
// 		expectedCode int
// 		expectedMsg  string
// 	}{
// 		{
// 			name: "User not registered in library",
// 			requestBody: models.RequestInput{
// 				ISBN:        "978-3-16-148410-0",
// 				LibraryId:   21, // Non-existent library ID
// 				RequestType: "borrow",
// 			},
// 			expectedCode: http.StatusBadRequest,
// 			expectedMsg:  "User is not authorized to request book user must be register in library",
// 		},
// 		{
// 			name: "Library does not exist",
// 			requestBody: models.RequestInput{
// 				ISBN:        "978-3-16-148410-0",
// 				LibraryId:   9999, // Non-existent library ID
// 				RequestType: "borrow",
// 			},
// 			expectedCode: http.StatusBadRequest,
// 			expectedMsg:  "Library is not present",
// 		},
// 		{
// 			name: "Book not in library",
// 			requestBody: models.RequestInput{
// 				ISBN:        "978-3-16-148411-0",
// 				LibraryId:   20, // Valid library but no book exists
// 				RequestType: "borrow",
// 			},
// 			expectedCode: http.StatusBadRequest,
// 			expectedMsg:  "Book is not present in library",
// 		},
// 		{
// 			name: "Successful borrow request",
// 			requestBody: models.RequestInput{
// 				ISBN:        "978-3-16-148410-0",
// 				LibraryId:   20,
// 				RequestType: "borrow",
// 			},
// 			expectedCode: http.StatusOK,
// 			expectedMsg:  "Request is successfully placed",
// 		},
// 		{
// 			name: "Duplicate borrow request",
// 			requestBody: models.RequestInput{
// 				ISBN:        "978-3-16-148411-1",
// 				LibraryId:   20,
// 				RequestType: "borrow",
// 			},
// 			expectedCode: http.StatusBadRequest,
// 			expectedMsg:  "Book is already requested",
// 		},
// 		{
// 			name: "Return request without borrowing",
// 			requestBody: models.RequestInput{
// 				ISBN:        "978-3-16-148412-1", // Different ISBN
// 				LibraryId:   20,
// 				RequestType: "return",
// 			},
// 			expectedCode: http.StatusBadRequest,
// 			expectedMsg:  "You have not taken any book",
// 		},
// 		{
// 			name: "Successful return request",
// 			requestBody: models.RequestInput{
// 				ISBN:        "978-3-16-148410-0",
// 				LibraryId:   20,
// 				RequestType: "return",
// 			},
// 			expectedCode: http.StatusOK,
// 			expectedMsg:  "Returned Book Successfully",
// 		},
// 	}

// 	// Run test cases
// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			jsonValue, _ := json.Marshal(tc.requestBody)
// 			req, _ := http.NewRequest("POST", "/request-event", bytes.NewBuffer(jsonValue))
// 			req.Header.Set("Content-Type", "application/json")
// 			w := httptest.NewRecorder()
// 			r.ServeHTTP(w, req)

// 			assert.Equal(t, tc.expectedCode, w.Code)

// 			var responseBody map[string]interface{}
// 			json.Unmarshal(w.Body.Bytes(), &responseBody)
// 			assert.Equal(t, tc.expectedMsg, responseBody["message"])
// 		})
// 	}
// }
