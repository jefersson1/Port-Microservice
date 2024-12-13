package event

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type emitter struct {
	connection *amqp.Connection
}

// push event to queue
func (e *emitter) Push(event string, severity string) error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	log.Println("pushing to channel...")

	err = channel.Publish(
		"auth_topic",
		severity,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func NewEventEmitter(conn *amqp.Connection) (emitter, error) {
	em := emitter{
		connection: conn,
	}

	// set it up
	channel, err := em.connection.Channel()
	if err != nil {
		return emitter{}, err
	}
	defer channel.Close()

	err = channel.ExchangeDeclare(
		"auth_topic",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return emitter{}, err
	}

	return em, nil
}
