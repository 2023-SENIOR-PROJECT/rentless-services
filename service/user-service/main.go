package main

import (
	"fmt"
	user_database "rentless-services/internal/infrastructure/user_database"
	controllers "rentless-services/internal/infrastructure/user_database/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	userDB := user_database.ConnectDatabase()

	// router.POST("/users", controllers.CreateUser)
	// router.GET("/users", controllers.GetAllUser)
	// router.GET("/users/:id", controllers.GetOneUser)
	// router.PUT("/users/:id", controllers.UpdateUser)
	// router.DELETE("/users/:id", controllers.DeleteUser)

	router.POST("/users", func(c *gin.Context) {
		controllers.CreateUser2(c, userDB)
	})
	router.GET("/users", func(c *gin.Context) {
		controllers.GetAllUser2(c, userDB)
	})
	router.GET("/users/:id", func(c *gin.Context) {
		controllers.GetOneUser2(c, userDB)
	})
	router.PUT("/users/:id", func(c *gin.Context) {
		controllers.UpdateUser2(c, userDB)
	})
	router.DELETE("/users/:id", func(c *gin.Context) {
		controllers.DeleteUser2(c, userDB)
	})

	fmt.Println("running in localhost:8080")
	router.Run("localhost:8080")
}
