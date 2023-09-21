// >> Unfinished code << 

package main

import (
	"github.com/gin-gonic/gin"
	"rentless-services/internal/infrastructure/user-database"
)

func main() {
	router := gin.Default()

	user_database.ConnectDatabase()

	// router.POST("/users", controllers.CreateUser)
	// router.GET("/users", controllers.GetUsers)
	// router.PUT("/users/:id", controllers.UpdateUser)
	// router.DELETE("/users/:id", controllers.DeleteUser)

	router.Run("localhost:8080")
}
