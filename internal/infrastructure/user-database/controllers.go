// >> Unfinished code << 

package user_database

import (
	"net/http"
	"rentless-services/internal/infrastructure/user-database/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	dsn := "user:password@tcp(127.0.0.1:3306)/database_name?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	DB.AutoMigrate(&models.User{})
}

type CreateUserInput struct {
	Firstname string `json:"firstname" binding:"required"`
	Lastname  string `json:"lastname" binding:"required"`
}

func CreateUser(c *gin.Context) {
	var input CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{Firstname: input.Firstname, Lastname: input.Lastname}
	DB.Create(&user)

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func GetUsers(c *gin.Context) {
	var users []models.User
	DB.Find(&users)

	c.JSON(http.StatusOK, gin.H{"data": users})
}

type UpdateUserInput struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Age       int    `json:"age"`
}

func UpdateUser(c *gin.Context) {
	var user models.User
	if err := DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "record not found"})
		return
	}

	var input UpdateUserInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser := models.User{Firstname: input.Firstname, Lastname: input.Lastname, Age: input.Age}

	DB.Model(&user).Updates(&updatedUser)
	c.JSON(http.StatusOK, gin.H{"data": user})
}

func DeleteUser(c *gin.Context) {
	var user models.User
	if err := DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "record not found"})
		return
	}

	DB.Delete(&user)
	c.JSON(http.StatusOK, gin.H{"data": "success"})
}
