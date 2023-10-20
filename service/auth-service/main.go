package main

import (
	"context"
	database "rentless-services/internal/infrastructure/product_database/mongo"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	database.ConnectMongoDB()
	r := gin.Default()
	r.POST("/auth/register", register)
	r.POST("/auth/login", login)
	r.GET("/auth/logout", logout)
	r.GET("/auth/validate", validateToken)
	r.Run(":8082")
}

type RegisterRequest struct {
	Email     string `json:"email" bson:"email"`
	Password1 string `json:"password1" bson:"password1"`
	Password2 string `json:"password2" bson:"password2"`
}
type LoginRequest struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type UserAuthSturct struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	Token    string `json:"token" bson:"token"`
}

func register(c *gin.Context) {
	var req RegisterRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "bad request",
		})
		return
	}
	if req.Password1 != req.Password2 {
		c.JSON(400, gin.H{
			"message": "passwords do not match",
		})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": req.Email,
	})
	tokenString, err := token.SignedString([]byte("secret"))
	var user UserAuthSturct
	user.Email = req.Email
	user.Password = req.Password1
	user.Token = tokenString
	database.DB.Collection("auth").InsertOne(context.Background(), user)
	c.JSON(200, gin.H{
		"message": "success",
	})
}

func login(c *gin.Context) {
	var req LoginRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "bad request",
		})
		return
	}
	var user UserAuthSturct
	err = database.DB.Collection("auth").FindOne(context.Background(), bson.M{"email": req.Email, "password": req.Password}).Decode(&user)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "invalid credentials",
			"error":   err.Error(),
		})
		return
	}
	c.SetCookie("token", user.Token, 3600, "/", "localhost", false, true)
	c.JSON(200, gin.H{"message": "success"})
}

func logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.JSON(200, gin.H{"message": "success"})
}

func validateToken(c *gin.Context) {
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.JSON(400, gin.H{"message": "invalid token"})
		return
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		c.JSON(400, gin.H{"message": "invalid token"})
		return
	}
	if !token.Valid {
		c.JSON(400, gin.H{"message": "invalid token"})
		return
	}
	c.JSON(200, gin.H{"message": "success"})
}
