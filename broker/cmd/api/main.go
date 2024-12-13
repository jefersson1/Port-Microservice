package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const port = "80"

func main() {
	log.Println("starting rabbitmq server...")
	conn, err := connectToRabbitMQ()
	if err != nil {
		log.Fatalf("could not connect to rabbitmq: %v", err)
	}
	defer conn.Close()

	// start broker server
	srv := newServer(conn).Router
	log.Println("starting broker server...")
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), srv)

	if errors.Is(err, http.ErrServerClosed) {
		log.Println("broker server closed")
	} else if err != nil {
		log.Println("error starting broker server:", err)
		os.Exit(1)
	}

}

func connectToRabbitMQ() (*amqp.Connection, error) {
	count := 0

	for {
		conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672")
		if err != nil {
			fmt.Println("rabbitmq not yet ready...")
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
