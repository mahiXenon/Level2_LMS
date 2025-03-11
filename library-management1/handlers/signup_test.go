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

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/create-user", CreateUser)
	return r
}

var testCases = []struct {
	name         string
	input        models.AuthInput
	expectedCode int
}{
	{
		name: "Successful user creation",
		input: models.AuthInput{
			Name:          "John Doe",
			Email:         "john.doe@example.com",
			Password:      "Password@123",
			ContactNumber: "1234567890",
		},
		expectedCode: http.StatusOK,
	},
	{
		name: "Duplicate email",
		input: models.AuthInput{
			Name:          "Jane Doe",
			Email:         "jane.doe@example.com",
			Password:      "Password@123",
			ContactNumber: "1122334455",
		},
		expectedCode: http.StatusBadRequest,
	},
	{
		name: "Invalid name (contains number)",
		input: models.AuthInput{
			Name:          "User123",
			Email:         "user123@example.com",
			Password:      "Password@123",
			ContactNumber: "2233445566",
		},
		expectedCode: http.StatusBadRequest,
	},
	{
		name: "Invalid email format",
		input: models.AuthInput{
			Name:          "John Doe",
			Email:         "invalid-email",
			Password:      "Password@123",
			ContactNumber: "1234567890",
		},
		expectedCode: http.StatusBadRequest,
	},
	{
		name: "Invalid contact number",
		input: models.AuthInput{
			Name:          "John Doe",
			Email:         "john.doe@example.com",
			Password:      "Password@123",
			ContactNumber: "123",
		},
		expectedCode: http.StatusBadRequest,
	},
	{
		name: "Invalid password",
		input: models.AuthInput{
			Name:          "John Doe",
			Email:         "john.doe@example.com",
			Password:      "pass",
			ContactNumber: "1234567890",
		},
		expectedCode: http.StatusBadRequest,
	},
}

func TestCreateUser(t *testing.T) {
	database.DB = database.SetDB() // Ensure database connection is set up properly
	r := SetupRouter()

	database.DB.Exec("DELETE FROM users") // Clear users before test

	password, _ := bcrypt.GenerateFromPassword([]byte("Password@123"), bcrypt.DefaultCost)
	database.DB.Create(&models.User{
		Name:          "Jane Doe",
		Email:         "jane.doe@example.com",
		Password:      string(password),
		ContactNumber: "0987654321",
		Role:          "user",
	})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonValue, _ := json.Marshal(tc.input)
			req, _ := http.NewRequest("POST", "/create-user", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, tc.expectedCode, w.Code)
		})
	}
}
