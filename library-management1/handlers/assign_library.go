package handlers

import (
	"fmt"
	"library-management1/database"
	"library-management1/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AssignAdmin(c *gin.Context) {
	user, err := c.Get("currentUser")
	if !err {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no user found"})
		return
	}

	var userData models.User
	userData = user.(models.User)
	if userData.Role == "Owner" {
		var admin models.Admin
		if err := c.ShouldBindJSON(&admin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var changeRole models.User
		database.DB.Where("id = ?", admin.ID).Find(&changeRole)
		if changeRole.ID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "User do not exist"})
			return
		}

		var library models.LibraryUser
		fmt.Println("name:", userData.Name)
		database.DB.Where("user_id = ?", userData.ID).Find(&library)
		fmt.Println("lib", library.LibraryId)

		// To check there should be one admin only
		type adminUsers []struct {
			UserId int `json:"id"`
		}

		var users adminUsers
		database.DB.Table("library_users as lb").
			Select("lb.user_id").
			Joins("join libraries as l on lb.library_id = l.id").
			Joins("join users as u on u.id = lb.user_id").
			Where("u.role = ?", "admin").
			Find(&users)

		if len(users) != 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "User is already Admin"})
			return
		}

		// to check already assigned a admin
		var checkAdmin models.User
		database.DB.Where("id = ?", admin.ID).Where("role = ?", "admin").Find(&checkAdmin)
		if checkAdmin.ID != 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "User is already Admin"})
			return
		}
		database.DB.Model(models.User{}).Where("id = ?", admin.ID).Update("Role", "admin")
		libraryData := models.LibraryUser{
			UserId:    admin.ID,
			LibraryId: library.LibraryId,
		}
		database.DB.Create(&libraryData)
		c.JSON(http.StatusOK, gin.H{"message": "Admin Assigned Successfully"})
		c.Set("currentUser", userData)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request Rejected"})
	}
}
