package logger

import (
	"context"
	"log"

	"go.uber.org/zap"
)

const (
	ctxKey loggerCtxKey = "logger_ctx_key"
)

type loggerCtxKey string

type Logger struct {
	*zap.Logger
}

func NewContext(ctx context.Context, logger interface{}) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if _, ok := logger.(*Logger); !ok {
		logger = newLogger()
	}

	return context.WithValue(ctx, ctxKey, logger)
}

func FromContext(ctx context.Context) *Logger {
	if logger, ok := ctx.Value(ctxKey).(*Logger); ok {
		return logger
	}

	return newLogger()
}

func newLogger() *Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Panic("error while init logger: %s ", err)
	}

	logger.Info("logger created")
	return &Logger{logger}
}
