package main

import (
	"fmt"
	"library-management1/database"
	"library-management1/handlers"
	"library-management1/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB()
	fmt.Println("Hello!")
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "You are on the home page"})
	})

	router.POST("/getdata", handlers.Test)
	router.DELETE("/del", handlers.Del)
	router.POST("/auth/signup", handlers.CreateUser)
	router.POST("/auth/login", handlers.Login)
	router.GET("/book/:search", middlewares.CheckAuth(), handlers.SearchBook)
	router.GET("/book/all", middlewares.CheckAuth(), handlers.GetAllBook)
	router.GET("/user/profile", middlewares.CheckAuth(), handlers.GetUserProfile)
	router.POST("/owner/create-library", middlewares.CheckAuth(), handlers.CreateLibrary)
	router.POST("/owner/assign-admin", middlewares.CheckAuth(), handlers.AssignAdmin)
	router.POST("/user/register", middlewares.CheckAuth(), handlers.Register)
	router.POST("/admin/insert-book", middlewares.CheckAuth(), handlers.InsertBook)
	router.POST("/admin/update-book", middlewares.CheckAuth(), handlers.UpadateBookCopies)
	router.POST("/user/request", middlewares.CheckAuth(), handlers.RequestEvent)
	router.GET("/admin/see-request", middlewares.CheckAuth(), handlers.SeeRequest)
	router.POST("/admin/handle-request", middlewares.CheckAuth(), handlers.HandleRequest)

	router.Run("localhost:8000")
}
