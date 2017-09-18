package logger

import (
	"context"
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	ctxKey loggerCtxKey = "logger_ctx_key"
)

type loggerCtxKey string

type ILogger interface {
	Info(string, ...zapcore.Field)
	Error(string, ...zapcore.Field)
	Debug(string, ...zapcore.Field)
	Panic(string, ...zapcore.Field)
}

// Logger is a logger service
type Logger struct {
	logger *zap.Logger
}

// NewContext stores logger in context
func NewContext(ctx context.Context, logger interface{}) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if _, ok := logger.(*Logger); !ok {
		logger = newLogger()
	}

	return context.WithValue(ctx, ctxKey, logger)
}

// FromContext returns logger form context
func FromContext(ctx context.Context) *Logger {
	if logger, ok := ctx.Value(ctxKey).(*Logger); ok {
		return logger
	}

	return newLogger()
}

func newLogger() *Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Panicf("error while init logger: %s ", err)
	}

	logger.Info("logger created")
	return &Logger{
		logger: logger,
	}
}

// Info is info level log
func (l *Logger) Info(msg string, fields ...zapcore.Field) {
	l.logger.Info(msg, fields...)
}

// Error is a error level log
func (l *Logger) Error(msg string, fields ...zapcore.Field) {
	l.logger.Error(msg, fields...)
}

// Debug is a debug level log
func (l *Logger) Debug(msg string, fields ...zapcore.Field) {
	l.logger.Debug(msg, fields...)
}

// Panic is a panic level log
func (l *Logger) Panic(msg string, fields ...zapcore.Field) {
	l.logger.Panic(msg, fields...)
}
