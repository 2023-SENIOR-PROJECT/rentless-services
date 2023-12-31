package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"review-consumer/db"
	"review-consumer/models"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

type Response struct {
	Data    Data   `json:"data"`
	Message string `json:"message"`
}

type Data struct {
	Email  string `json:"email"`
	UserID uint   `json:"user_id"`
}

const (
	RabbitMQURL  = "amqps://qntphlli:2dL5SE3y0b43BU_1xJHQtcCXO5BibvTz@armadillo.rmq.cloudamqp.com/qntphlli"
	QueueName    = "review_queue"
	ExchangeName = "review_exchange"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func GetConnection(url string) *amqp.Connection {
	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")
	return conn
}

func GetChannel(connection *amqp.Connection) *amqp.Channel {
	ch, err := connection.Channel()
	failOnError(err, "Failed to open a channel")
	return ch
}

func Publish(message, rountingKey string, channel *amqp.Channel) {
	err := channel.Publish(
		"",          // exchange
		rountingKey, // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	failOnError(err, "Failed to publish a message")
}

var CHANNEL = GetChannel(GetConnection(RabbitMQURL))

func validate(token string) (uint, error) {
	// Create a new request for validation
	validateRequest, err := http.NewRequest(http.MethodGet, "http://user-service:8080/auth/validate", nil)
	if err != nil {
		fmt.Println("Error creating validation request:", err)
		return 0, err
	}

	// Set the token in the request header
	validateRequest.Header.Set("Cookie", "token="+token)
	validateRequest.Header.Set("Authorization", "Bearer "+token)

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

func GetAllReviews(c *gin.Context, db *db.ReviewDB) {
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

func GetAllReviewsByProductID(c *gin.Context, db *db.ReviewDB) {
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

func GetOneReview(c *gin.Context, db *db.ReviewDB) {
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

func CreateReview(c *gin.Context, db *db.ReviewDB) {
	// token, err := c.Cookie("token")
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(401, gin.H{"message": "no token"})
		return
	}
	tokenString = tokenString[7:]
	authorID, err := validate(tokenString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("authorID: ", authorID)
	productID := c.Param("productID")
	var review models.Review
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review.AuthorID = authorID
	review.ProductID = productID

	body, err := json.Marshal(review)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = CHANNEL.Publish(
		"",        // Exchange
		QueueName, // Routing key
		false,     // Mandatory
		false,     // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(" [x] Sent", string(body))

	// err = db.CreateRecord(authorID, productID, review)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }
	c.JSON(http.StatusCreated, gin.H{"message": "Review Created"})
}
