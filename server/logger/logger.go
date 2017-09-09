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

// Logger is a logger service
type Logger struct {
	*zap.Logger

	prefix string
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
		Logger: logger,
	}
}

func (l *Logger) Prefix(prefix string) *Logger {
	l.prefix = prefix

	return l
}
