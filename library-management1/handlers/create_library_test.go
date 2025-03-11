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

func SetupLibraryRouter() *gin.Engine {
	r := gin.Default()
	r.Use(MockMiddleware()) // Middleware to set user context
	r.POST("/create-library", CreateLibrary)
	return r
}

// MockMiddleware to simulate authentication and set user context
func MockMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		database.DB.Where("email = ?", "john.doe@example.com").First(&user)
		if user.ID != 0 {
			c.Set("currentUser", user)
		}
		c.Next()
	}
}

func TestCreateLibrary(t *testing.T) {
	database.DB = database.SetDB()
	r := SetupLibraryRouter()

	// Clear libraries and users before testing
	database.DB.Exec("DELETE FROM libraries")
	database.DB.Exec("DELETE FROM users")

	// Create a test user
	user := models.User{
		Name:          "John Doe",
		Email:         "john.doe@example.com",
		Password:      "Password@123",
		ContactNumber: "1234567890",
		Role:          "user",
	}
	database.DB.Create(&user)

	testCases := []struct {
		name         string
		input        models.AuthLibrary
		expectedCode int
		expectedRole string
		setUser      bool
	}{
		{
			name: "Successful library creation",
			input: models.AuthLibrary{
				Name: "Central Library",
			},
			expectedCode: http.StatusOK,
			expectedRole: "Owner",
			setUser:      true,
		},
		{
			name: "Duplicate library name",
			input: models.AuthLibrary{
				Name: "Central Library",
			},
			expectedCode: http.StatusBadRequest,
			expectedRole: "user",
			setUser:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonValue, _ := json.Marshal(tc.input)
			req, _ := http.NewRequest("POST", "/create-library", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)

			// if tc.setUser {
			// 	var updatedUser models.User
			// 	database.DB.Where("email = ?", user.Email).Find(&updatedUser)
			// 	assert.Equal(t, tc.expectedRole, updatedUser.Role)
			// }
		})
	}
}
