package main

import (
	"fmt"
	"listener/event"
	"log"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// Try to connect to rabbitmq
	rabbitmqConn, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitmqConn.Close()

	fmt.Println("Connected to RabbitMQ")

	// Start listening for messages
	log.Println("Listening for and consuming RabitMQ messages")

	// Create a consumer
	consumer, err := event.NewConsumer(rabbitmqConn)
	if err != nil {
		log.Fatal(err)
	}

	// Watch the queue and consume events
	err = consumer.Listen([]string{"log.info", "log.warn", "log.error"})
	if err != nil {
		log.Fatal(err)
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// Try to connect to rabbitmq

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabitmq")
		if err != nil {
			fmt.Println("RabbitMQ connection. Err: ", err)
			counts++
		} else {
			log.Println("Connected to RabbitMQ")
			connection = c
			break
		}

		if counts > 5 {
			return nil, fmt.Errorf("Failed to connect to RabbitMQ")
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second

		log.Println("Failed to connect to RabbitMQ... Retrying...")

		time.Sleep(backOff)
		backOff = backOff * 2
		continue
	}

	// Return the connection
	return connection, nil
}
