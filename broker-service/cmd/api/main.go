package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "80"

type Config struct {
	RabbitMQ *amqp.Connection
}

func main() {
	// Try to connect to rabbitmq
	rabbitmqConn, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitmqConn.Close()

	fmt.Println("Connected to RabbitMQ")

	app := Config{
		RabbitMQ: rabbitmqConn,
	}

	log.Printf("Starting broker service on port %s\n", webPort)

	// define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start the server
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
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
