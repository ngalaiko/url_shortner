package logger

import "context"

// NewContext stores logger in context
func NewContext(ctx context.Context, logger interface{}) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if _, ok := logger.(ILogger); !ok {
		logger = newLogger()
	}

	return context.WithValue(ctx, ctxKey, logger)
}

// FromContext returns logger form context
func FromContext(ctx context.Context) ILogger {
	if logger, ok := ctx.Value(ctxKey).(ILogger); ok {
		return logger
	}

	return newLogger()
}
