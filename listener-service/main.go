package main

import (
	"fmt"
	"listener-service/event"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// connect to RabbitMQ
	conn, err := connect()
	if err != nil {
		log.Panic(err)
		os.Exit(1)
	}
	defer conn.Close()

	// start listening for messages
	log.Println("Listening for and consuming messages...")

	// create consumer
	consumer, err := event.NewConsumer(conn, "logs_topic")
	if err != nil {
		log.Panic(err)
		os.Exit(1)
	}

	// watch the queue and consume events
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Panic(err)
		os.Exit(1)
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	for {
		c, err := amqp.Dial("amqp://guest:guest@localhost")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			log.Println("Connected to RabbitMQ!")
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println("Failed to connect after 5 attempts")
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		fmt.Println("Backing off...")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
