package main

import (
	"log"
	"net/http"
	database "rentless-services/internal/infrastructure/rental_database/mongo"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

// CreateRental creates a new rental
func CreateRental(c *gin.Context) {
	var product Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := database.InsertOne(product)
	c.JSON(http.StatusCreated, result)
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

	c.JSON(http.StatusNoContent, nil)
}

func parseObjectID(id string, c *gin.Context) (primitive.ObjectID, error) {
	pr, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	return pr, err
}

type Product struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	OrderID   string             `json:"orderId" bson:"orderId"`
	ProductID string             `json:"productId" bson:"productId"`
	Quantity  int                `json:"quantity" bson:"quantity"`
	Amount    float64            `json:"amount" bson:"amount"`
}
