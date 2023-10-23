package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	review_database "rentless-services/internal/infrastructure/review_database"

	models "rentless-services/internal/infrastructure/review_database/models"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Data    Data   `json:"data"`
	Message string `json:"message"`
}

type Data struct {
	Email  string `json:"email"`
	UserID uint   `json:"user_id"`
}

func validate(token string) (uint, error) {
	// Create a new request for validation
	validateRequest, err := http.NewRequest(http.MethodGet, "http://localhost:8080/auth/validate", nil)
	if err != nil {
		fmt.Println("Error creating validation request:", err)
		return 0, err
	}

	// Set the token in the request header
	validateRequest.Header.Set("Cookie", "token="+token)

	// Send the validation request
	validateResponse, err := http.DefaultClient.Do(validateRequest)
	if err != nil {
		fmt.Println("Error sending validation request:", err)
		return 0, err
	}
	defer validateResponse.Body.Close()

	var validateBody Response
	if err := json.NewDecoder(validateResponse.Body).Decode(&validateBody); err != nil {
		fmt.Println("err: ", err)
	}

	fmt.Println(token)
	fmt.Println(validateBody.Data.UserID)

	return validateBody.Data.UserID, nil
}

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
	token, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	authorID, err := validate(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	productID := c.Param("productID")
	var review models.Review
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = db.CreateRecord(authorID, productID, review)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Review Created"})
}
