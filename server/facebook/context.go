package facebook

import (
	"context"
)

type ctxKey string

const (
	apiCtxKey ctxKey = "context_key_for_facebook_api"
)

// NewContext stores facebook api in context
func NewContext(ctx context.Context, facebook interface{}) context.Context {

	if ctx == nil {
		ctx = context.Background()
	}

	if facebook == nil {
		facebook = newApi(ctx)
	}

	return context.WithValue(ctx, apiCtxKey, facebook)
}

// FromContext returns facebook api from context
func FromContext(ctx context.Context) *Api {
	if api, ok := ctx.Value(apiCtxKey).(*Api); ok {
		return api
	}

	return newApi(ctx)
}
