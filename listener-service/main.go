package main

import (
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"listener/event"
	"log"
	"math"
	"os"
	"time"
)

func main() {
	// try to connect to rabbit
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	// start listening to message
	log.Println("Listening for and consuming RabbitMQ messages...")

	// create consume
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}

	//watch the queue and consume events
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
	}
}

func connect() (*amqp091.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp091.Connection
	var rabbitURL = os.Getenv("RABBIT_URL")

	for {
		c, err := amqp091.Dial(rabbitURL)
		if err != nil {
			fmt.Print("RabbitMQ not yet ready...")
			counts++
		} else {
			connection = c
			fmt.Println()
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		fmt.Printf("Backing off for %d seconds...\n", int(math.Pow(float64(counts), 2)))
		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
