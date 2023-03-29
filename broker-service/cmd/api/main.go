package main

import (
	"errors"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
	"os"
	"time"
)

const webPort = ":80"

type Config struct {
	Rabbit *amqp091.Connection
}

func main() {
	rabbitConn, err := connectToRabbit()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	app := Config{
		Rabbit: rabbitConn,
	}

	log.Printf("Starting broker service on port %s\n", webPort)

	svr := &http.Server{
		Addr:    fmt.Sprintf("%s", webPort),
		Handler: app.routes(),
	}

	err = svr.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connectToRabbit() (*amqp091.Connection, error) {
	var rabbitConn *amqp091.Connection
	var count int64
	var rabbitURL = os.Getenv("RABBIT_URL")

	for {
		connection, err := amqp091.Dial(rabbitURL)
		if err != nil {
			fmt.Println("Rabbit not ready...")
			count++
		} else {
			fmt.Println()
			rabbitConn = connection
			break
		}

		if count > 15 {
			fmt.Println(err)
			return nil, errors.New("cannot connect to rabbit")
		}

		fmt.Println("Backing off for 2 seconds...")
		time.Sleep(2 * time.Second)
		continue
	}

	fmt.Println("Connected to rabbit")
	return rabbitConn, nil
}
