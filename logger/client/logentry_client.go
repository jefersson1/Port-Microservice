package client

type LogEntryClient interface {
	InsertLogEntry(logEntry LogEntry) error
}
