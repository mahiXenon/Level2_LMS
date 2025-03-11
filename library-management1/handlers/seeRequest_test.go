package handlers

import (
	"encoding/json"
	// "fmt"
	"library-management1/database"
	"library-management1/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Setup router for testing
func SetupSeeRequest() *gin.Engine {
	r := gin.Default()

	r.Use(MockMiddlewareSeeRequest()) // Middleware to set user context
	r.GET("/see-request", SeeRequest)
	return r
}

// // MockMiddlewareForEventRouter simulates authentication and sets user context
func MockMiddlewareSeeRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		database.DB.Where("email = ?", "user@example.com").First(&user)
		if user.ID != 0 {
			c.Set("currentUser", user)
		}
		c.Next()
	}
}

func TestSeeRequest(t *testing.T) {
	r := SetupSeeRequest()
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
	database.DB.Create(&user)

	user1 := models.User{

		Name:          "Test User",
		Email:         "user1@example.com",
		Password:      "Password@123",
		ContactNumber: "1234567890",
		Role:          "admin",
	}
	database.DB.Create(&user1)

	request := models.RequestEvent{
		ISBN:        "978-3-16-148410-0",
		UserId:      6,
		LibraryId:   2,
		RequestType: "borrow",
	}
	database.DB.Create(&request)
	// // gin.SetMode(gin.TestMode)
	// tests := []struct {
	// 	name string
	// 	// requestBody  models.RequestInput
	// 	expectedCode int
	// 	expectedMsg  string
	// }{
	// 	{
	// 		name:         "Non -admin User",
	// 		expectedCode: http.StatusBadRequest,
	// 		expectedMsg:  "Only admin can see all request",
	// 	},
	// }
	libraryUser := models.LibraryUser{
		UserId:    user.ID,
		LibraryId: 10,
	}
	database.DB.Create(&libraryUser)

	t.Run("No request found", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/see-request", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var responseBody map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.Equal(t, []interface{}{}, responseBody["data"]) // Expect empty list
	})

	t.Run("Request found", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/see-request", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// 	var responseBody map[string]interface{}
		// 	json.Unmarshal(w.Body.Bytes(), &responseBody)
		// 	request, exists := responseBody["data"].([]interface{})
		// 	assert.True(t, exists, "Data should be a list of books")
		// 	assert.Len(t, request, 1, "There should be 1 request")
	})

	// tests := []struct {
	// 	name           string
	// 	currentUser    models.User
	// 	expectedStatus int
	// 	expectedBody   string
	// }{
	// 	{
	// 		name:           "Non-admin user",
	// 		currentUser:    models.User{Role: "user"},
	// 		expectedStatus: http.StatusBadRequest,
	// 		expectedBody:   "Only admin can see all request",
	// 	},
	// 	{
	// 		name:           "Admin user with no requests",
	// 		currentUser:    models.User{Role: "admin", ID: 1},
	// 		expectedStatus: http.StatusOK,
	// 		expectedBody:   `"data":[]`,
	// 	},
	// 	{
	// 		name:           "Admin user with requests",
	// 		currentUser:    models.User{Role: "admin", ID: 1},
	// 		expectedStatus: http.StatusOK,
	// 		expectedBody:   `"data":[`,
	// 	},
	// }

	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		router := gin.Default()
	// 		router.GET("/seeRequest", func(c *gin.Context) {
	// 			c.Set("currentUser", tt.currentUser)
	// 			SeeRequest(c)
	// 		})

	// 		req, _ := http.NewRequest(http.MethodGet, "/seeRequest", nil)
	// 		resp := httptest.NewRecorder()
	// 		router.ServeHTTP(resp, req)

	//			assert.Equal(t, tt.expectedStatus, resp.Code)
	//			assert.Contains(t, resp.Body.String(), tt.expectedBody)
	//		})
	//	}
	// for _, tc := range tests {
	// 	t.Run(tc.name, func(t *testing.T) {
	// 		// jsonValue, _ := json.Marshal(tc.requestBody)
	// 		req, _ := http.NewRequest("POST", "/see-request", nil)
	// 		req.Header.Set("Content-Type", "application/json")
	// 		w := httptest.NewRecorder()
	// 		r.ServeHTTP(w, req)

	// 		assert.Equal(t, tc.expectedCode, w.Code)

	// 		// var responseBody map[string]interface{}
	// 		var responseBody map[string]interface{}
	// 		// json.Unmarshal(w.Body.Bytes(), &responseBody)
	// 		// assert.Equal(t, []interface{}{}, responseBody["data"])
	// 		books, exists := responseBody["data"].([]interface{})
	// 		assert.True(t, exists, "Data should be a list of books")
	// 		assert.Len(t, books, 1, "There should be 1 book in response")
	// 	})
	// }
}
