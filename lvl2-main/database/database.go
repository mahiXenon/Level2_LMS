package database

import (
	"fmt"
	"library-management1/models"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB

func InitDB() {
	// Connect to the PostgreSQL database
	var err error
	DB, err = gorm.Open(
		"postgres",
		"host=localhost port=5433 user=postgres "+
			"dbname=postgres password=postgres sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	//	DB = db
	DB.AutoMigrate(&models.User{}, &models.Library{}, &models.LibraryUser{}, &models.BookInventory{}, &models.RequestEvent{}, &models.IssueBook{})

	fmt.Println("Connected!")
	// defer db.Close()
}
