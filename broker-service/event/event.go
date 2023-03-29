package event

import "github.com/rabbitmq/amqp091-go"

func getExchangeName() string {
	return "logs_topic"
}

func declareExchange(ch *amqp091.Channel) error {
	return ch.ExchangeDeclare(
		getExchangeName(),
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
}
