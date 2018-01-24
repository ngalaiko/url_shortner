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
		facebook = newAPI(ctx)
	}

	return context.WithValue(ctx, apiCtxKey, facebook)
}

// FromContext returns facebook api from context
func FromContext(ctx context.Context) *API {
	if api, ok := ctx.Value(apiCtxKey).(*API); ok {
		return api
	}

	return newAPI(ctx)
}
