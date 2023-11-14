package main

import (
	"encoding/json"
	"fmt"
	"log"
	"review-consumer/db"
	"review-consumer/models"

	"github.com/streadway/amqp"
)

const (
	RabbitMQURL  = "amqps://qntphlli:2dL5SE3y0b43BU_1xJHQtcCXO5BibvTz@armadillo.rmq.cloudamqp.com/qntphlli"
	QueueName    = "review_queue"
	ExchangeName = "review_exchange"
	RoutingKey   = "review"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial(RabbitMQURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		QueueName, // Queue name
		false,     // Durable
		false,     // Delete when unused
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		failOnError(err, "Failed to declare a queue")
	}

	msgs, err := ch.Consume(
		q.Name, // Queue
		"",     // Consumer
		true,   // Auto-Ack
		false,  // Exclusive
		false,  // No-local
		false,  // No-Wait
		nil,    // Args
	)
	if err != nil {
		failOnError(err, "Failed to register a consumer")
	}

	forever := make(chan bool)

	// Connect Database
	reviewDB := db.ConnectDatabase()
	fmt.Println(reviewDB)

	go func() {
		for d := range msgs {
			// Recive message algor
			fmt.Printf("Received a message: %s\n", d.Body)
			var review models.Review
			err := json.Unmarshal(d.Body, &review)
			if err != nil {
				panic(err)
			}
			// fmt.Println("id: ", review.ID)
			// fmt.Println("createdAt: ", review.CreatedAt)
			// fmt.Println("authorID: ", review.AuthorID)
			// fmt.Println("productID: ", review.ProductID)
			// fmt.Println("rate: ", review.Rate)
			// fmt.Println("comment: ", review.Comment)
			err = reviewDB.CreateRecord(review.AuthorID, review.ProductID, review)
			if err != nil {
				fmt.Println("err: ", err)
			}
		}
	}()

	fmt.Println(" [*] Waiting for messages. To exit, press CTRL+C")
	<-forever
}
