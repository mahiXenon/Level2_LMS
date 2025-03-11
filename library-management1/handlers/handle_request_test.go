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

// func SetupHandleRouter() *gin.Engine {
// 	r := gin.Default()
// 	r.Use(MockMiddlewareForEventRouter()) // Middleware to set user context
// 	r.POST("/handle_request", HandleRequest)
// 	return r
// }

// // MockMiddlewareForEventRouter simulates authentication and sets user context
// func MockMiddlewareForHandleRouter() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var user models.User
// 		database.DB.Where("email = ?", "user@example.com").First(&user)
// 		if user.ID != 0 {
// 			c.Set("currentUser", user)
// 		}
// 		c.Next()
// 	}
// }

// func TestHandleRequest(t *testing.T) {
// 	r := SetupHandleRouter()
// 	database.DB = database.SetDB()

// 	database.DB.Exec("DELETE FROM libraries")
// 	database.DB.Exec("DELETE FROM users")
// 	database.DB.Exec("DELETE FROM library_users")
// 	database.DB.Exec("DELETE FROM book_inventories")
// 	database.DB.Exec("DELETE FROM request_events")
// 	user := models.User{

// 		Name:          "Test User",
// 		Email:         "user@example.com",
// 		Password:      "Password@123",
// 		ContactNumber: "1234567890",
// 		Role:          "admin",
// 	}
// 	database.DB.Create(&user) // Insert the user

// 	request := models.RequestEvent{
// 		ISBN:        "1234567891",
// 		UserId:      user.ID,
// 		LibraryId:   3,
// 		RequestType: "borrow",
// 	}
// 	database.DB.Create(&request)

// 	request1 := models.RequestEvent{
// 		ISBN:        "1234567890",
// 		UserId:      user.ID,
// 		LibraryId:   3,
// 		RequestType: "borrow",
// 	}
// 	database.DB.Create(&request1)

// 	request2 := models.RequestEvent{
// 		ISBN:        "1234567890",
// 		UserId:      user.ID,
// 		LibraryId:   3,
// 		RequestType: "return",
// 	}
// 	database.DB.Create(&request2)

// 	// request2.RequestType = "return"
// 	// database.DB.Model(&models.RequestEvent{}).Where("id = ?", request2.ID).Update("request_type = ?", request2.RequestType)

// 	book := models.BookInventory{
// 		ISBN:            "1234567890",
// 		LibraryId:       3,
// 		Title:           "Book Title",
// 		Author:          "Author",
// 		Publisher:       "Publisher",
// 		Version:         "1.0",
// 		TotalCopies:     10,
// 		AvailableCopies: 10,
// 	}
// 	database.DB.Create(&book)
// 	tests := []struct {
// 		name         string
// 		user         models.User
// 		requestBody  models.ReaderRequest
// 		expectedCode int
// 		expectedBody string
// 	}{
// 		// {
// 		// 	name: "Unauthorized user",
// 		// 	user: models.User{Role: "user"},
// 		// 	requestBody: models.ReaderRequest{
// 		// 		ID: 1,
// 		// 	},
// 		// 	expectedCode: http.StatusBadRequest,
// 		// 	expectedBody: "Yoar are not authorized to approve request",
// 		// },
// 		{
// 			name: "Request not found",
// 			// user: models.User{Role: "admin"},
// 			requestBody: models.ReaderRequest{
// 				ID: 999,
// 			},
// 			expectedCode: http.StatusBadRequest,
// 			expectedBody: "Request not found",
// 		},
// 		{
// 			name: "Book not available",
// 			// user: models.User{Role: "admin"},
// 			requestBody: models.ReaderRequest{
// 				ID: request.ID,
// 			},
// 			expectedCode: http.StatusBadRequest,
// 			expectedBody: "Book is not available",
// 		},
// 		{
// 			name: "Borrow request approved",
// 			// user: models.User{Role: "admin"},
// 			requestBody: models.ReaderRequest{
// 				ID: request1.ID,
// 			},
// 			expectedCode: http.StatusOK,
// 			expectedBody: "Request is approved",
// 		},
// 		{
// 			name: "Return request approved",
// 			// user: models.User{Role: "admin"},
// 			requestBody: models.ReaderRequest{
// 				ID: request2.ID,
// 			},
// 			expectedCode: http.StatusOK,
// 			expectedBody: "Book is submitted to liberary",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			jsonValue, _ := json.Marshal(tt.requestBody)

// 			// body, _ := json.Marshal(tt.requestBody)
// 			req, _ := http.NewRequest(http.MethodPost, "/handle_request", bytes.NewBuffer(jsonValue))
// 			req.Header.Set("Content-Type", "application/json")

// 			w := httptest.NewRecorder()
// 			r.ServeHTTP(w, req)
// 			var responseBody map[string]interface{}
// 			json.Unmarshal(w.Body.Bytes(), &responseBody)
// 			assert.Equal(t, tt.expectedCode, w.Code)
// 			// assert.Equal(t, tt.expectedBody, responseBody["message"])
// 		})
// 	}
// }
