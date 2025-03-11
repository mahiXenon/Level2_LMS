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

func SetupRegisterRouter() *gin.Engine {
	r := gin.Default()
	r.Use(MockMiddlewareForRegistration())
	r.POST("/register", Register)
	return r
}

func MockMiddlewareForRegistration() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		database.DB.Where("email = ?", "user@example.com").First(&user)
		if user.ID != 0 {
			c.Set("currentUser", user)
		}
		c.Next()
	}
}

func TestRegister(t *testing.T) {
	database.DB = database.SetDB()
	r := SetupRegisterRouter()

	// Clear database before testing
	database.DB.Exec("DELETE FROM library_users")
	database.DB.Exec("DELETE FROM libraries")
	database.DB.Exec("DELETE FROM users")

	// Create a regular user
	user := models.User{
		Name:          "Test User",
		Email:         "user@example.com",
		Password:      "Password@123",
		ContactNumber: "1234567890",
		Role:          "user",
	}
	database.DB.Create(&user)

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

	testCases := []struct {
		name         string
		input        models.RegisterLibrary
		expectedCode int
		expectedBody map[string]interface{}
		mockUser     models.User
	}{
		{
			name: "Successful registration",
			input: models.RegisterLibrary{
				Name: "Main Library",
			},
			expectedCode: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "Sucessfully Registered",
			},
			mockUser: user,
		},
		{
			name: "Library does not exist",
			input: models.RegisterLibrary{
				Name: "Unknown Library",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"message": "No library found",
			},
			mockUser: user,
		},
		// {
		// 	name: "Admin cannot register in a library again",
		// 	input: models.RegisterLibrary{
		// 		Name: "Main Library",
		// 	},
		// 	expectedCode: http.StatusBadRequest,
		// 	expectedBody: map[string]interface{}{
		// 		"message": "You are not allowed to register in other library",
		// 	},
		// 	mockUser: admin,
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonValue, _ := json.Marshal(tc.input)
			req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")

			// Mock user in request context
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)

			var responseBody map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(t, err, "Response should be valid JSON")

			for key, expectedValue := range tc.expectedBody {
				assert.Equal(t, expectedValue, responseBody[key], "Mismatch in response for key: "+key)
			}
		})
	}
}
