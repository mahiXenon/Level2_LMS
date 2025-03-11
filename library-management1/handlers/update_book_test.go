package handlers

// package handlers

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

func SetupUpdateRouter() *gin.Engine {
	r := gin.Default()
	r.Use(MockMiddlewareUpdateRouter()) // Middleware to set user context
	r.POST("/updateBook", UpadateBookCopies)
	return r
}

// MockMiddlewareForEventRouter simulates authentication and sets user context
func MockMiddlewareUpdateRouter() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		database.DB.Where("email = ?", "user@example.com").First(&user)
		if user.ID != 0 {
			c.Set("currentUser", user)
		}
		c.Next()
	}
}
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		database.DB.Where("email = ?", "user1@example.com").First(&user)
		if user.ID != 0 {
			c.Set("currentUser", user)
		}
		c.Next()
	}
}
func TestUpdateBookCopies(t *testing.T) {
	// gin.SetMode(gin.TestMode)
	r := SetupUpdateRouter()
	database.DB = database.SetDB()
	database.DB.Exec("DELETE FROM libraries")
	database.DB.Exec("DELETE FROM users")
	database.DB.Exec("DELETE FROM library_users")
	database.DB.Exec("DELETE FROM book_inventories")
	database.DB.Exec("DELETE FROM request_events")
	user := models.User{

		Name:          "Test User",
		Email:         "user@example.com",
		Password:      "Password@123",
		ContactNumber: "1234567890",
		Role:          "admin",
	}
	database.DB.Create(&user) // Insert the user

	library := models.LibraryUser{
		UserId:    user.ID,
		LibraryId: 3,
	}
	database.DB.Create(&library)

	book := models.BookInventory{
		ISBN:            "1234567890",
		LibraryId:       3,
		Title:           "Book Title",
		Author:          "Author",
		Publisher:       "Publisher",
		Version:         "1.0",
		TotalCopies:     10,
		AvailableCopies: 10,
	}

	database.DB.Create(&book)

	tests := []struct {
		name           string
		currentUser    models.User
		updateBook     models.UpdateBookDetails
		expectedStatus int
		expectedBody   string
	}{
		// {
		// 	name:           "Invalid JSON",
		// 	// currentUser:    models.User{},
		// 	updateBook:     models.UpdateBookDetails{},
		// 	expectedStatus: http.StatusBadRequest,
		// 	expectedBody:   `{"error":"invalid character 'i' looking for beginning of object key string"}`,
		// },

		{
			name:           "Book not present in library",
			currentUser:    models.User{Role: "admin"},
			updateBook:     models.UpdateBookDetails{ISBN: "1234567891", ADD: 5, DecreaseCount: 2},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"Book is not present in library"}`,
		},
		{
			name:           "Decrease count more than total copies",
			currentUser:    models.User{Role: "admin"},
			updateBook:     models.UpdateBookDetails{ISBN: "1234567890", ADD: 0, DecreaseCount: 13},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"Decrease count is more than total copies"}`,
		},
		{
			name: "Successful update",
			// currentUser:    models.User{Role: "admin"},
			updateBook:     models.UpdateBookDetails{ISBN: "1234567890", ADD: 5, DecreaseCount: 2},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"data":{}}`, // Adjust this based on actual response
		},
	}
	t.Run("Only admin can update book", func(t *testing.T) {

		updateBook := models.UpdateBookDetails{
			ISBN: "1234567891",
			ADD:  5, DecreaseCount: 2,
		}

		q := gin.Default()
		database.DB = database.SetDB()
		// r.Use(MockMiddlewareUpdateRouter())
		// Middleware to set user context
		q.Use(Middleware())
		q.POST("/updateBook", UpadateBookCopies)
		user1 := models.User{

			Name:          "Test User",
			Email:         "user1@example.com",
			Password:      "Password@123",
			ContactNumber: "1238567890",
			Role:          "user",
		}
		database.DB.Create(&user1)
		body, _ := json.Marshal(updateBook)
		req, _ := http.NewRequest("POST", "/updateBook", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		q.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		// var response map[string]interface{}
		// json.Unmarshal(w.Body.Bytes(), &response)
		// assert.Equal(t, authInput.Email, response["data"].(map[string]interface{})["email"])
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonValue, _ := json.Marshal(tt.updateBook)
			// router := gin.Default()
			// router.PUT("/updateBook", func(c *gin.Context) {
			// 	c.Set("currentUser", tt.currentUser)
			// 	UpadateBookCopies(c)
			// })

			// var body []byte
			// if tt.name == "Invalid JSON" {
			// 	body = []byte(`{invalid json}`)
			// } else {
			// 	body, _ = json.Marshal(tt.updateBook)
			// }

			req, _ := http.NewRequest(http.MethodPost, "/updateBook", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			r.ServeHTTP(resp, req)
			var responseBody map[string]interface{}
			json.Unmarshal(resp.Body.Bytes(), &responseBody)
			assert.Equal(t, tt.expectedStatus, resp.Code)
			// assert.Equal(t, tt.expectedBody, responseBody["message"])
			// assert.Contains(t, resp.Body.String(), tt.expectedBody)
		})
	}
}
