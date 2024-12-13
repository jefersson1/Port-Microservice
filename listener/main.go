package main

import (
	"log"
	"os"
	"time"

	event "github.com/jateen67/listener/wabbit"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	rabbitConn, err := connectToRabbitMQ()
	if err != nil {
		log.Fatalf("could not connect to rabbitmq: %v", err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	log.Println("listening for and consuming rabbitmq messages...")

	// consumer
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		log.Fatalf("could not create new consumer: %v", err)
	}

	// queue watch
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Fatalf("could not listen to queue: %v", err)
	}
}

func connectToRabbitMQ() (*amqp.Connection, error) {
	count := 0

	for {
		conn, err := amqp.Dial(os.Getenv("RABBITMQ_CONNECTION_STRING"))
		if err != nil {
			log.Println("rabbitmq not yet ready...")
			count++
		} else {
			log.Println("connected to rabbitmq successfully")
			return conn, nil
		}

		if count > 10 {
			log.Println(err)
			return nil, err
		}

		log.Println("retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
	}
}
