package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	// "gorm.io/driver/sqlite"

	"library-management1/database"
	"library-management1/models"
)

// func setupTestDB() {
// 	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
// 	db.AutoMigrate(&models.Library{})
// 	database.DB = db
// }

func TestGetAllLibrary(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database.DB = database.SetDB()

	// Insert test data
	database.DB.Create(&models.Library{Name: "Test Library 1"})
	database.DB.Create(&models.Library{Name: "Test Library 2"})

	router := gin.Default()
	router.GET("/libraries", GetAllLibrary)

	req, _ := http.NewRequest(http.MethodGet, "/libraries", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Test Library 1")
	assert.Contains(t, w.Body.String(), "Test Library 2")
}
