package session

import "context"

const (
	ctxKey = "session_ctx_key"
)

// NewContext places session service to context
func NewContext(ctx context.Context, session interface{}) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if _, ok := session.(ISession); !ok {
		session = newSession(ctx)
	}

	return context.WithValue(ctx, ctxKey, session)
}

// FromContext returns session service form context
func FromContext(ctx context.Context) ISession {
	if session, ok := ctx.Value(ctxKey).(ISession); ok {
		return session
	}
	return newSession(ctx)
}
