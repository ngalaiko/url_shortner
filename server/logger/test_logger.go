package logger

import (
	"log"

	"go.uber.org/zap"
)

// NewTestLogger returns debug level logger
func NewTestLogger() *Logger {
	l, err := zap.NewDevelopment()
	if err != nil {
		log.Panicf("error while init logger: %s ", err)
	}

	return &Logger{l}
}
