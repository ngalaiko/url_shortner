package cache

import "context"

// NewContext stores cache in context
func NewContext(ctx context.Context, cache interface{}) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if _, ok := cache.(ICache); !ok {
		cache = newCache(ctx)
	}

	return context.WithValue(ctx, ctxKey, cache)
}

// FromContext returns cache from context
func FromContext(ctx context.Context) ICache {
	if cache, ok := ctx.Value(ctxKey).(ICache); ok {
		return cache
	}

	return newCache(ctx)
}
