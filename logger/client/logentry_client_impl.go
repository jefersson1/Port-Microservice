package client

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

type logEntryClientImpl struct {
	Client *mongo.Client
}

func NewLogEntryClientImpl(c *mongo.Client) *logEntryClientImpl {
	return &logEntryClientImpl{Client: c}
}

func (l *logEntryClientImpl) InsertLogEntry(logEntry LogEntry) error {
	coll := l.Client.Database("logs_db").Collection("logs")
	doc := LogEntry{Name: logEntry.Name, Data: logEntry.Data, CreatedAt: logEntry.CreatedAt}

	_, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		log.Println("error inserting log entry:", err)
		return err
	}

	return nil
}
