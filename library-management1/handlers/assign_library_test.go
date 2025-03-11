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

func SetupAssignAdminRouter() *gin.Engine {
	r := gin.Default()
	r.Use(MockMiddlewareForAdmin()) // Middleware to set user context
	r.POST("/assign-admin", AssignAdmin)
	return r
}

// MockMiddleware to simulate authentication and set user context
func MockMiddlewareForAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		database.DB.Where("email = ?", "owner@example.com").First(&user)
		if user.ID != 0 {
			c.Set("currentUser", user)
		}
		c.Next()
	}
}

func TestAssignAdmin(t *testing.T) {
	database.DB = database.SetDB()
	r := SetupAssignAdminRouter()

	// Clear database before testing
	database.DB.Exec("DELETE FROM library_users")
	database.DB.Exec("DELETE FROM libraries")
	database.DB.Exec("DELETE FROM users")

	// Create an owner user
	owner := models.User{
		Name:          "Owner User",
		Email:         "owner@example.com",
		Password:      "Password@123",
		ContactNumber: "1234567890",
		Role:          "Owner",
	}
	database.DB.Create(&owner)

	// Create a normal user to be assigned as admin
	normalUser := models.User{
		Name:          "Normal User",
		Email:         "normal@example.com",
		Password:      "Password@123",
		ContactNumber: "0987654321",
		Role:          "user",
	}
	database.DB.Create(&normalUser)

	// Create a library and link the owner
	library := models.Library{
		Name: "Main Library",
	}
	database.DB.Create(&library)

	libraryUser := models.LibraryUser{
		UserId:    owner.ID,
		LibraryId: library.ID,
	}
	database.DB.Create(&libraryUser)

	testCases := []struct {
		name         string
		input        models.Admin
		expectedCode int
		expectedBody map[string]interface{}
	}{
		{
			name: "Successful admin assignment",
			input: models.Admin{
				ID: normalUser.ID,
			},
			expectedCode: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "Admin Assigned Successfully",
			},
		},
		{
			name: "Admin already assigned",
			input: models.Admin{
				ID: normalUser.ID,
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"message": "User is already Admin",
			},
		},
		{
			name: "User does not exist",
			input: models.Admin{
				ID: 9999, // Non-existing user ID
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"message": "User do not exist",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonValue, _ := json.Marshal(tc.input)
			req, _ := http.NewRequest("POST", "/assign-admin", bytes.NewBuffer(jsonValue))
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
				assert.Equal(t, expectedValue, responseBody[key], "Mismatch in response for key: "+key)
			}
		})
	}
}
