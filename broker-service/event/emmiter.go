package event

import (
	"context"
	"github.com/rabbitmq/amqp091-go"
	"log"
)

type Emitter struct {
	connection *amqp091.Connection
}

func (e *Emitter) setup() error {
	channel, err := e.connection.Channel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()
	return declareExchange(channel)
}

func (e *Emitter) Push(event string, severity string) error {
	channel, err := e.connection.Channel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()

	log.Println("Pushing to", getExchangeName())

	err = channel.PublishWithContext(
		context.Background(),
		getExchangeName(),
		severity,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Sending message: %s -> %s", event, getExchangeName())
	return nil
}

func NewEventEmitter(conn *amqp091.Connection) (Emitter, error) {
	emitter := Emitter{
		connection: conn,
	}

	err := emitter.setup()
	if err != nil {
		return Emitter{}, err
	}

	return emitter, nil
}
