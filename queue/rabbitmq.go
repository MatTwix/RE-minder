package queue

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/MatTwix/RE-minder/config"
	"github.com/rabbitmq/amqp091-go"
)

var (
	conn    *amqp091.Connection
	channel *amqp091.Channel
)

func Connect() {
	cfg := config.LoadConfig()
	conn, err := amqp091.Dial(cfg.RabbitMQUrl)
	if err != nil {
		log.Fatalf("Error connecting to RabbitMQ: %v", err)
	}

	channel, err = conn.Channel()
	if err != nil {
		log.Fatalf("Error opening a channel: %v", err)
	}

	log.Println("Successfully connected to RabbitMQ")
}

func Publish(queueName string, body []byte) error {
	q, err := channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return errors.New("Error declaring queue: " + err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = channel.PublishWithContext(
		ctx,
		"",
		q.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp091.Persistent,
		},
	)

	return err
}

func Close() {
	if channel != nil {
		channel.Close()
	}
	if conn != nil {
		conn.Close()
	}
}
