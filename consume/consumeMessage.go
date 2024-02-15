package main

import (
	"log"

	"github.com/streadway/amqp"
)

func main() {
	// URL of the RabbitMQ server
	brokerURL := "amqps://yourusername:yourpassword@your-amq-broker-url:port"
	// Name of the queue to consume
	queueName := "myQueue"

	if err := consumeMessages(brokerURL, queueName); err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}
}

func consumeMessages(brokerURL, queueName string) error {
	// Connect to RabbitMQ server
	conn, err := amqp.Dial(brokerURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Open a channel
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Declare the queue (to ensure it exists)
	q, err := ch.QueueDeclare(
		queueName,
		false, // Durable
		false, // Delete when not used
		false, // Exclusive
		false, // No wait
		nil,   // Arguments
	)
	if err != nil {
		return err
	}

	// Receive messages
	msgs, err := ch.Consume(
		q.Name, // Queue name
		"",     // Consumer tag
		false,  // Auto-acknowledge set to false, manual confirmation of message handling
		false,  // Exclusive
		false,  // No local
		false,  // No wait
		nil,    // No arguments
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	// Use a goroutine to listen for messages
	go func() {
		for d := range msgs {
			// Process the message
			log.Printf("Received a message: %s", d.Body)

			// Manually acknowledge the message after successful processing
			if err := d.Ack(false); err != nil {
				log.Printf("Error acknowledging message: %s", err)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	// Block the main goroutine and wait for a signal (e.g., wait for user input or others)
	// This is to allow the consumers goroutine to continue to work without being closed
	<-forever

	return nil
}
