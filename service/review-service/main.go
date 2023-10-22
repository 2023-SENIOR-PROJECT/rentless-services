package main

import (
	"fmt"
	review_database "rentless-services/internal/infrastructure/review_database"
	controllers "rentless-services/internal/infrastructure/review_database/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	reviewDB := review_database.ConnectDatabase()
	fmt.Println(reviewDB)

	// Get one review by reviewID
	router.GET("/review/:reviewID", func(c *gin.Context) {
		controllers.GetOneReview(c, reviewDB)
	})
	//Get all reviews
	router.GET("/reviews", func(c *gin.Context) {
		controllers.GetAllReviews(c, reviewDB)
	})
	// Get reviews by productID
	router.GET("/reviews/:productID", func(c *gin.Context) {
		controllers.GetAllReviewsByProductID(c, reviewDB)
	})
	// Create review
	router.POST("/reviews/:productID", func(c *gin.Context) {
		controllers.CreateReview(c, reviewDB)
	})

	defer reviewDB.DB.Close()
	fmt.Println("running in localhost:8080")
	router.Run("localhost:8080")
}
