package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/jateen67/logger/client"
	logger "github.com/jateen67/logger/protos"
	"google.golang.org/grpc"
)

const port = "50001"

func main() {
	// start mongo
	mongoClient, err := client.ConnectToClient()
	if err != nil {
		log.Fatalf("could not connect to mongo: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// close connection
	defer func() {
		err = mongoClient.Disconnect(ctx)
		if err != nil {
			panic(err)
		}
	}()

	logEntryClient := client.NewLogEntryClientImpl(mongoClient)
	// start logger server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	logger.RegisterLoggerServiceServer(s, &server{loggerClient: logEntryClient})

	log.Println("starting grpc server...")
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
