package controllers

import (
	"authservice/db"
	"authservice/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// New Version
// New Get One User API
func GetOneUser2(c *gin.Context, db *db.UserDB) {
	id := c.Param("id")
	if !db.UserExists(id) {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	user, err := db.GetOneUser(id)
	if !db.UserExists(id) {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// New Get All User API
func GetAllUser2(c *gin.Context, db *db.UserDB) {
	users, err := db.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(users) == 0 {
		c.JSON(http.StatusNoContent, gin.H{"message": "No users in database"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// New Create User API Not done
func CreateUser2(c *gin.Context, db *db.UserDB) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result_user, err := db.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result_user)
}

// New Update User Not done
func UpdateUser2(c *gin.Context, db *db.UserDB) {
	id := c.Param("id")
	var user models.User
	fmt.Println("Before Context to JSON")
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(user)
	if !db.UserExists(id) {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	fmt.Println("After Exists")
	user, err := db.UpdateUser(id, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// New Delete User done
func DeleteUser2(c *gin.Context, db *db.UserDB) {
	id := c.Param("id")
	if !db.UserExists(id) {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	err := db.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

// auth service
func Register(c *gin.Context, userDB *db.UserDB) {
	var req models.RegisterRequest
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
	userInfo := models.User{
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Age:       req.Age,
	}
	res, err := userDB.CreateUser(userInfo)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "bad request",
			"error":   err.Error(),
		})
		return
	}
	fmt.Println(res)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":   req.Email,
		"user_id": res.ID,
	})
	tokenString, err := token.SignedString([]byte("secret"))
	var user models.UserAuthSturct
	user.Email = req.Email
	user.Pwd = req.Password1
	user.Token = tokenString
	// database.DB.Collection("auth").InsertOne(context.Background(), user)

	err = userDB.CreateNewUser(user, res.ID)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "bad request",
			"error":   err.Error(),
		})
		return
	}

	registerResponse := models.LoginResponse{
		Firstname: res.Firstname,
		Email:     req.Email,
		Token:     tokenString,
		User_id:   res.ID,
	}

	c.JSON(200, gin.H{
		"message": "success",
		"data":    registerResponse,
	})
}

func Login(c *gin.Context, userDB *db.UserDB) {
	var req models.LoginRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "bad request",
		})
		return
	}
	var user models.UserAuthSturct

	// err = database.DB.Collection("auth").FindOne(context.Background(), bson.M{"email": req.Email, "password": req.Pwd}).Decode(&user)
	user, err = userDB.FindUserAccount(req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "invalid credentials",
			"error":   err.Error(),
		})
		return
	}

	uidstr := strconv.FormatUint(uint64(user.User_id), 10)
	fmt.Println(uidstr)
	userInfo, err := userDB.GetOneUser(uidstr)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "invalid credentials",
			"error":   err.Error(),
		})
		return
	}
	loginResponse := models.LoginResponse{
		Firstname: userInfo.Firstname,
		Email:     user.Email,
		Token:     user.Token,
		User_id:   user.User_id,
	}

	c.SetCookie("token", user.Token, 3600, "/", "localhost", false, true)
	c.JSON(200, gin.H{"message": "success", "data": loginResponse})
}

func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.JSON(200, gin.H{"message": "success"})
}

func ValidateToken(c *gin.Context) {
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.JSON(401, gin.H{"message": "invalid token"})
		return
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		c.JSON(401, gin.H{"message": "invalid token"})
		return
	}
	fmt.Println(token.Claims)
	if !token.Valid {
		c.JSON(401, gin.H{"message": "invalid token"})
		return
	}
	c.JSON(200, gin.H{"message": "success", "data": token.Claims})
}
