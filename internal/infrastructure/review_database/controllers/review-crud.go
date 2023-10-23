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
	avgRateAndCount, err := db.GetAvgRateAndCountAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(reviews) == 0 {
		c.JSON(http.StatusNoContent, gin.H{"message": "No reviews in database"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"avg_rate": avgRateAndCount.AvgRate, "number_reviews": avgRateAndCount.NumberReview, "reviews": reviews})
}

func GetAllReviewsByProductID(c *gin.Context, db *review_database.ReviewDB) {
	productID := c.Param("productID")
	reviews, err := db.GetRecordsByProductID(productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	avgRateAndCount, err := db.GetAvgRateAndCountByProductID(productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(reviews) == 0 {
		c.JSON(http.StatusNoContent, gin.H{"message": "No reviews in database"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"avg_rate": avgRateAndCount.AvgRate, "number_reviews": avgRateAndCount.NumberReview, "reviews": reviews})
}

func GetOneReview(c *gin.Context, db *review_database.ReviewDB) {
	reviewID := c.Param("reviewID")
	if !db.ReviewExists(reviewID) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}
	review, err := db.GetOneRecord(reviewID) //If review is not exist will return []
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"review": review})
}

func CreateReview(c *gin.Context, db *review_database.ReviewDB) {
	authorID := "1"
	productID := c.Param("productID")
	var review models.Review
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := db.CreateRecord(authorID, productID, review)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Review Created"})
}
