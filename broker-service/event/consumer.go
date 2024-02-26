package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

func NewConsumer(conn *amqp.Connection) (*Consumer, error) {

	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		return nil, err
	}

	return &Consumer{
		conn: conn,
	}, nil
}

func (c *Consumer) setup() error {

	channel, err := c.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (c *Consumer) Listen(topics []string) error {
	channel, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	q, err := declareRandomQueue(channel)
	if err != nil {
		return err
	}

	for _, topic := range topics {
		err = channel.QueueBind(
			q.Name,       // queue name
			topic,        // routing key
			"logs_topic", // exchange
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	msgs, err := channel.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			var payload = Payload{}
			_ = json.Unmarshal(d.Body, &payload)

			go handlePayload(payload)
		}
	}()

	fmt.Printf("Waiting for message [Exchange, Queue] [logs_topic, %s]\n", q.Name)
	<-forever

	return nil
}

func handlePayload(payload Payload) {
	fmt.Printf("Received a message: %s\n", payload)

	switch payload.Name {
	case "log", "event":
		// Log event
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}

	case "auth":
		// Authenticate event

	default:
		log.Println("Unknown event")
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}

	}
}

func logEvent(payload Payload) error {
	jsonData, _ := json.MarshalIndent(payload, "", "\t")
	request, err := http.NewRequest("POST", "http://logger-service/log", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}
