package logger

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	ctxKey loggerCtxKey = "logger_ctx_key"
)

type loggerCtxKey string

// ILogger is an interface for logger
type ILogger interface {
	Info(string, ...zapcore.Field)
	Error(string, ...zapcore.Field)
	Debug(string, ...zapcore.Field)
	Panic(string, ...zapcore.Field)
}

// Logger is a logger service
type logger struct {
	logger *zap.Logger
}

func newLogger() *logger {
	l, err := zap.NewProduction()
	if err != nil {
		log.Panicf("error while init logger: %s ", err)
	}

	l.Info("logger created")
	return &logger{
		logger: l,
	}
}

// Info is info level log
func (l *logger) Info(msg string, fields ...zapcore.Field) {
	l.logger.Info(msg, fields...)
}

// Error is a error level log
func (l *logger) Error(msg string, fields ...zapcore.Field) {
	l.logger.Error(msg, fields...)
}

// Debug is a debug level log
func (l *logger) Debug(msg string, fields ...zapcore.Field) {
	l.logger.Debug(msg, fields...)
}

// Panic is a panic level log
func (l *logger) Panic(msg string, fields ...zapcore.Field) {
	l.logger.Panic(msg, fields...)
}
