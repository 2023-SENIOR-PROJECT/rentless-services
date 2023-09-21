package main

import (
	"fmt"
	user_database "rentless-services/internal/infrastructure/user_database"
	controllers "rentless-services/internal/infrastructure/user_database/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	user_database.ConnectDatabase()

	router.POST("/users", controllers.CreateUser)
	router.GET("/users", controllers.GetAllUser)
	router.GET("/users/:id", controllers.GetOneUser)
	router.PUT("/users/:id", controllers.UpdateUser)
	router.DELETE("/users/:id", controllers.DeleteUser)

	fmt.Println("running in localhost:8080")
	router.Run("localhost:8080")
}
