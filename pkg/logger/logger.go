package logger

import (
	"log"
)

type Logger interface {
	Info(message string)
	Error(message string)
}

type SimpleLogger struct{}

func NewLogger() Logger {
	return &SimpleLogger{}
}

func (l *SimpleLogger) Info(message string) {
	log.Println("INFO:", message)
}

func (l *SimpleLogger) Error(message string) {
	log.Println("ERROR:", message)
}
