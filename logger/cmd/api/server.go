package main

import (
	"context"
	"log"
	"time"

	"github.com/jateen67/logger/client"
	logger "github.com/jateen67/logger/protos"
)

type server struct {
	logger.UnimplementedLoggerServiceServer
	loggerClient client.LogEntryClient
}

func (s *server) LogActivity(ctx context.Context, req *logger.LogRequest) (*logger.LogResponse, error) {
	doc := client.LogEntry{
		Name:      req.Name,
		Data:      req.Data,
		CreatedAt: time.Now(),
	}

	err := s.loggerClient.InsertLogEntry(doc)
	if err != nil {
		log.Fatalf("failed to insert log into database: %v", err)
		return nil, err
	}

	res := &logger.LogResponse{
		Error:   false,
		Message: "Succesfully logged activity!",
		LogEntry: &logger.LogEntry{
			Name:      doc.Name,
			Data:      doc.Data,
			CreatedAt: doc.CreatedAt.Format(time.RFC3339),
		},
	}

	log.Println("logger service: successful log")
	return res, nil
}
