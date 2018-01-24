package token

import "context"

const (
	ctxKey tokenCtxKey = "user_token_ctx_key"
)

type tokenCtxKey string

// NewContext places user_token to context
func NewContext(ctx context.Context, tokens interface{}) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if _, ok := tokens.(*Service); !ok {
		tokens = newTokens(ctx)
	}

	return context.WithValue(ctx, ctxKey, tokens)
}

// FromContext returns links form context
func FromContext(ctx context.Context) *Service {
	if tokens, ok := ctx.Value(ctxKey).(*Service); ok {
		return tokens
	}

	return newTokens(ctx)
}
