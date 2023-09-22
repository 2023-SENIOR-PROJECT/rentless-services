package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	user_database "rentless-services/internal/infrastructure/user_database"
	models "rentless-services/internal/infrastructure/user_database/models"

	"github.com/gin-gonic/gin"
)

var db = user_database.GetDB()

// Get one user
func GetOneUser(c *gin.Context) {
	id := c.Param("id")

	query := "SELECT id, firstname, lastname, age, created_at, updated_at FROM users WHERE id = ?"
	row := db.QueryRow(query, id)
	var user models.User

	err := row.Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Age, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, user)
}

// GetAllUser retrieves all users
func GetAllUser(c *gin.Context) {
	// Define a SQL query to select all users
	query := "SELECT id, firstname, lastname, age, created_at, updated_at FROM users"

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Initialize a slice to hold the result
	var users []models.User

	// Loop through the result set and scan each row into a User struct
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Age, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		users = append(users, user)
	}

	if len(users) == 0 {
		c.JSON(http.StatusNoContent, gin.H{"message": "No users found"})
		return
	}

	c.JSON(http.StatusOK, users)
	defer rows.Close()
}

// Create a new user
func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Define a SQL query to insert a new user
	query := "INSERT INTO users (firstname, lastname, age, created_at, updated_at) VALUES (?, ?, ?, NOW(), NOW())"

	// Execute the query
	result, err := db.Exec(query, user.Firstname, user.Lastname, user.Age)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userID, _ := result.LastInsertId()
	user.ID = uint(userID)
	c.JSON(http.StatusCreated, user)
}

// Update an existing user by ID
func UpdateUser(c *gin.Context) {
	// Get the user ID from the URL parameter
	id := c.Param("id")

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !userExists(id) {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	query := "UPDATE users SET firstname = ?, lastname = ?, age = ?, updated_at = NOW() WHERE id = ?"

	_, err := db.Exec(query, user.Firstname, user.Lastname, user.Age, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Delete a user by ID
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if !userExists(id) {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	query := "DELETE FROM users WHERE id = ?"

	_, err := db.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

func userExists(id string) bool {
	query := "SELECT id FROM users WHERE id = ?"
	var userID uint
	err := db.QueryRow(query, id).Scan(&userID)
	return err != sql.ErrNoRows
}

// New Get One User API
func GetOneUser2(c *gin.Context, db *user_database.UserDB) {
	id := c.Param("id")
	user, err := db.GetOneUser(id)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found", "message": "No Row Error"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Get One Function Error"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// New Get All User API
func GetAllUser2(c *gin.Context, db *user_database.UserDB) {
	users, err := db.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Get All Function Error"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Context can not be JSON Error"})
		return
	}
	result_user, err := db.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Create Function Error"})
		return
	}
	c.JSON(http.StatusCreated, result_user)
}

// New Update User Not done
func UpdateUser2(c *gin.Context, db *user_database.UserDB) {
	id := c.Param("id")
	var user models.User
	fmt.Println("Before Conext to JSON")
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Context can not be JSON Error"})
		return
	}
	fmt.Println(user)
	if !userExists(id) {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	fmt.Println("After Exists")
	user, err := db.UpdateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Update Function Error"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// New Delete User done
func DeleteUser2(c *gin.Context, db *user_database.UserDB) {
	id := c.Param("id")
	if !userExists(id) {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	err := db.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Delete Function Error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
