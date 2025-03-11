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
	"golang.org/x/crypto/bcrypt"
)

func SetupLoginRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/login", Login)
	return r
}

func TestLogin(t *testing.T) {
	database.DB = database.SetDB()
	r := SetupLoginRouter()

	// Clear users before testing
	database.DB.Exec("DELETE FROM users")

	// Create a test user
	password, _ := bcrypt.GenerateFromPassword([]byte("Password@123"), bcrypt.DefaultCost)
	database.DB.Create(&models.User{
		Name:          "mahendra singh",
		Email:         "mahendra.singh@example.com",
		Password:      string(password),
		ContactNumber: "1234567890",
		Role:          "user",
	})

	testCases := []struct {
		name         string
		input        models.AuthLogin
		expectedCode int
	}{
		{
			name: "Successful login",
			input: models.AuthLogin{
				Email:    "mahendra.singh@example.com",
				Password: "Password@123",
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "User does not exist",
			input: models.AuthLogin{
				Email:    "nonexistent@example.com",
				Password: "Password@123",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Wrong password",
			input: models.AuthLogin{
				Email:    "mahendra.singh@example.com",
				Password: "WrongPassword",
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonValue, _ := json.Marshal(tc.input)
			req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, tc.expectedCode, w.Code)
		})
	}
}
