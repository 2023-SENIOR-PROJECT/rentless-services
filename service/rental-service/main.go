package main

import (
	"log"
	"net/http"
	database "rentless-services/internal/infrastructure/rental_database/mongo"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreateRentalRequest struct {
	ItemsPrice      float64                `json:"itemsPrice"`
	OrderItems      []OrderItem            `json:"orderItems"`
	PaymentMethod   string                 `json:"paymentMethod"`
	ShippingAddress map[string]interface{} `json:"shippingAddress"`
	ShippingPrice   float64                `json:"shippingPrice"`
	TaxPrice        float64                `json:"taxPrice"`
	TotalPrice      float64                `json:"totalPrice"`
}

type OrderItem struct {
	ID           string  `json:"_id"`
	Name         string  `json:"name"`
	CountInStock int     `json:"countInStock"`
	Image        string  `json:"image"`
	Price        float64 `json:"price"`
	Quantity     int     `json:"quantity"`
	Slug         string  `json:"slug"`
}

type Product struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	OrderID   string             `json:"orderId" bson:"orderId"`
	ProductID string             `json:"productId" bson:"productId"`
	Quantity  int                `json:"quantity" bson:"quantity"`
	Amount    float64            `json:"amount" bson:"amount"`
}

const (
	RabbitMQURL  = "amqp://guest:guest@localhost:5672/"
	QueueName    = "rental_queue"
	ExchangeName = "rental_exchange"
	RoutingKey   = "rental"
)

func main() {
	r := gin.Default()
	r.GET("/rentals", GetAllRentals)
	r.GET("/rentals/:id", GetRental)
	r.POST("/rentals", CreateRental)
	r.PUT("/rentals/:id", UpdateRental)
	r.DELETE("/rentals/:id", DeleteRental)

	_, err := database.ConnectMongoDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	port := ":8083"
	if err := r.Run(port); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}

func CreateRental(c *gin.Context) {
	var request CreateRentalRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a connection to RabbitMQ
	conn, err := amqp.Dial(RabbitMQURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to RabbitMQ"})
		return
	}
	defer conn.Close()

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open a channel"})
		return
	}
	defer ch.Close()

	// Declare a queue
	_, err = ch.QueueDeclare(QueueName, true, false, false, false, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to declare a queue"})
		return
	}

	// Publish each product as a message to RabbitMQ
	for _, orderItem := range request.OrderItems {
		// Create a JSON representation of the product
		productJSON := []byte(`{"productId": "` + orderItem.ID + `", "quantity": ` +
			strconv.Itoa(orderItem.Quantity) + `, "amount": ` +
			strconv.FormatFloat(orderItem.Price, 'f', -1, 64) + `}`)

		// Publish the product to RabbitMQ
		err = ch.Publish(ExchangeName, RoutingKey, false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        productJSON,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish a message to RabbitMQ"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Rental created successfully"})
}

// GetAllRentals retrieves all rentals
func GetAllRentals(c *gin.Context) {
	cursor, err := database.GetAllProduct()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(c.Request.Context())

	var products []Product
	if err := cursor.All(c.Request.Context(), &products); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetRental retrieves a single rental by ID
func GetRental(c *gin.Context) {
	id := c.Param("id")
	pr, err := parseObjectID(id, c)
	if err != nil {
		return
	}

	filter := bson.M{"_id": pr}
	singleResult := database.GetOneProduct(filter)

	if singleResult.Err() == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		return
	} else if singleResult.Err() != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": singleResult.Err().Error()})
		return
	}

	var product Product
	if err := singleResult.Decode(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// UpdateRental updates a rental by ID
func UpdateRental(c *gin.Context) {
	id := c.Param("id")
	pr, err := parseObjectID(id, c)
	if err != nil {
		return
	}

	var product Product
	if err := c.ShouldBindWith(&product, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := bson.M{"_id": pr}
	update := bson.M{"$set": product}

	result, err := database.UpdateOne(filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// DeleteRental deletes a rental by ID
func DeleteRental(c *gin.Context) {
	id := c.Param("id")
	pr, err := parseObjectID(id, c)
	if err != nil {
		return
	}

	filter := bson.M{"_id": pr}
	result, err := database.DeleteOne(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Product deleted successfully"})
}

func parseObjectID(id string, c *gin.Context) (primitive.ObjectID, error) {
	pr, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	return pr, err
}
