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

	// err := consumeMessages(brokerURL, queueName)
	// if err != nil {
	// 	log.Fatalf("Failed to consume messages: %v", err)
	// }
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

func consumeMessages(brokerURL, queueName string) error {
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

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			// 这里处理你的逻辑
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

	return nil
}
