package logger

import (
	"log"

	"go.uber.org/zap"
)

func NewTestLogger() *Logger {
	l, err := zap.NewDevelopment()
	if err != nil {
		log.Panicf("error while init logger: %s ", err)
	}

	return &Logger{l}
}
