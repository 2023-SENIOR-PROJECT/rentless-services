package main

import (
	"authservice/controllers"
	"authservice/db"
	"fmt"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	router := gin.Default()
	router.Use(CORSMiddleware())

	userDB := db.ConnectDatabase()

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
	router.POST("/auth/register", func(c *gin.Context) {
		controllers.Register(c, userDB)
	})
	router.POST("/auth/login", func(c *gin.Context) {
		controllers.Login(c, userDB)
	})
	router.GET("/auth/logout", controllers.Logout)
	router.GET("/auth/validate", controllers.ValidateToken)
	defer userDB.DB.Close()
	fmt.Println("running in localhost:8080")
	router.Run("localhost:8080")
}
