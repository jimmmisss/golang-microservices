package event

import "github.com/rabbitmq/amqp091-go"

func getExchangeName() string {
	return "logs_topic"
}

func declareRandomQueue(ch *amqp091.Channel) (amqp091.Queue, error) {
	return ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
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
