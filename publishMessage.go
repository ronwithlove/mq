package main

import (
	"log"

	"github.com/streadway/amqp"
)

func main() {
	//Amazon MQ Broker
	brokerURL := "amqps://yourusername:yourpassword@your-amq-broker-url:port"
	//The name of the queue created in Amazon MQ
	queueName := "myQueue"
	message := "Hello, this is a test message!"

	// Invoke the publishMessage function to send a message.
	err := publishMessage(brokerURL, queueName, message)
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	} else {
		log.Printf("Message published: %s", message)
	}
}

func publishMessage(brokerURL, queueName, message string) error {
	conn, err := amqp.Dial(brokerURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	return err
}
