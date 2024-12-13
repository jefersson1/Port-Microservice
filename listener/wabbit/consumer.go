package event

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type consumer struct {
	conn *amqp.Connection
}

type payload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewConsumer(conn *amqp.Connection) (consumer, error) {

	cons := consumer{
		conn: conn,
	}

	// set up the consumer by opening up a channel and declaring an exchange
	channel, err := cons.conn.Channel()
	if err != nil {
		return consumer{}, err
	}

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
		return consumer{}, err
	}

	return cons, nil
}

// listens to the queue for specific topics
func (consumer *consumer) Listen(topics []string) error {
	// go to our consumer channel and get things from it
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// we have our channel, now we need to get a random queue
	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// we have our channel and our queue now

	// go through our list of topics
	for _, s := range topics {
		// bind our channel to each of these topics
		err = ch.QueueBind(
			q.Name,
			s,
			"auth_topic",
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	// look for messages
	messages, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// declaring a new channel
	forever := make(chan bool)
	// will run in background
	go func() {
		for d := range messages {

			var p payload
			_ = json.Unmarshal(d.Body, &p)

			go handleLog(p)
		}
	}()

	log.Printf("waiting for message [exchange, queue] [auth_topic, %s]\n", q.Name)
	// keep the consumption going forever by making this blocking
	<-forever

	return nil
}

func handleLog(payload payload) {
	err := logEvent(payload)
	if err != nil {
		log.Println(err)
	}
}

func logEvent(entry payload) error {

	jsonData, _ := json.Marshal(entry)

	request, err := http.NewRequest("POST", "http://authentication/authentication", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// make sure we get the correct status code from the logger service
	if res.StatusCode != http.StatusOK {
		return err
	}

	return nil
}
