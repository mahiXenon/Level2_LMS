package database

import (
	"fmt"
	"library-management1/models"
	"log"

	"github.com/jinzhu/gorm"
)

// var Test_DB *gorm.DB

func SetDB() *gorm.DB {
	// var err error
	Test_DB, err := gorm.Open(
		"postgres",
		"host=localhost port=5433 user=postgres "+
			"dbname=test1 password=postgres sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	//	DB = db
	Test_DB.AutoMigrate(&models.User{}, &models.Library{}, &models.LibraryUser{}, &models.BookInventory{}, &models.RequestEvent{}, &models.IssueBook{})

	fmt.Println("Connected!")
	return Test_DB
}
