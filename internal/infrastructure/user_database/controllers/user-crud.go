package controllers

import (
	"fmt"
	"net/http"
	user_database "rentless-services/internal/infrastructure/user_database"
	models "rentless-services/internal/infrastructure/user_database/models"

	"github.com/gin-gonic/gin"
)

// New Version
// New Get One User API
func GetOneUser2(c *gin.Context, db *user_database.UserDB) {
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
func GetAllUser2(c *gin.Context, db *user_database.UserDB) {
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
func CreateUser2(c *gin.Context, db *user_database.UserDB) {
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
func UpdateUser2(c *gin.Context, db *user_database.UserDB) {
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
func DeleteUser2(c *gin.Context, db *user_database.UserDB) {
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
