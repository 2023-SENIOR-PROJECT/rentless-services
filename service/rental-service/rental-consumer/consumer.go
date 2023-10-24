package main

import (
	"encoding/json"
	"log"
	database "rentless-services/internal/infrastructure/rental_database/mongo"

	"github.com/streadway/amqp"
)

const (
	RabbitMQURL  = "amqp://guest:guest@localhost:5672/"
	QueueName    = "rental_queue"
	ExchangeName = "rental_exchange"
	RoutingKey   = "rental"
)

func main() {
	// Create a connection to RabbitMQ
	conn, err := amqp.Dial(RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	log.Printf("Connected to RabbitMQ")
	defer conn.Close()

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declare a queue
	_, err = ch.QueueDeclare(QueueName, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// Create a consumer for the queue
	msgs, err := ch.Consume(QueueName, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	// Consume and process messages
	for msg := range msgs {
		// Save to the database)
		processRental(msg.Body)

		log.Printf("Received a message: %s", msg.Body)
	}
}

func processRental(message []byte) {
	type Product struct {
		ProductID string  `json:"productId"`
		Quantity  int     `json:"quantity"`
		Amount    float64 `json:"amount"`
	}
	//unmarshal the message into a product struct
	var product Product
	err := json.Unmarshal(message, &product)
	if err != nil {
		log.Fatalf("Failed to unmarshal the message: %v", err)
	}
	//save the rental data to the database
	rental := database.InsertOne(product)
	if rental == nil {
		log.Fatalf("Failed to save rental data to the database: %v", err)
	}
}
