package controllers

import (
	"net/http"
	review_database "rentless-services/internal/infrastructure/review_database"

	models "rentless-services/internal/infrastructure/review_database/models"

	"github.com/gin-gonic/gin"
)

func GetAllReviews(c *gin.Context, db *review_database.ReviewDB) {
	reviews, err := db.GetAllRecords()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(reviews) == 0 {
		c.JSON(http.StatusNoContent, gin.H{"message": "No users in database"})
		return
	}
	c.JSON(http.StatusOK, reviews)
}

func GetAllReviewsByProductID(c *gin.Context, db *review_database.ReviewDB) {
	productID := c.Param("productID")
	reviews, err := db.GetRecordsByProductID(productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(reviews) == 0 {
		c.JSON(http.StatusNoContent, gin.H{"message": "No users in database"})
		return
	}
	c.JSON(http.StatusOK, reviews)
}

func GetOneReview(c *gin.Context, db *review_database.ReviewDB) {
	reviewID := c.Param("reviewID")
	review, err := db.GetOneRecord(reviewID) //If review is not exist will return []
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(review) == 0 {
		c.JSON(http.StatusNoContent, gin.H{"message": "No users in database"})
		return
	}
	c.JSON(http.StatusOK, review)
}

func CreateReview(c *gin.Context, db *review_database.ReviewDB) {
	var review models.Review
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// err := db.CreateRecord(review)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }
	c.JSON(http.StatusCreated, gin.H{"message": "Review Created", "body": review})
}
