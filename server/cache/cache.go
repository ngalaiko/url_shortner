package cache

import (
	"context"

	"golang.org/x/sync/syncmap"

	"github.com/ngalayko/url_shortner/server/logger"
	"go.uber.org/zap"
)

const (
	ctxKey cacheContextKey = "cache_ctx_key"
)

type cacheContextKey string

// Cache is a cache service
type Cache struct {
	logger *logger.Logger

	cacheMap *syncmap.Map
}

// NewContext stores cache in context
func NewContext(ctx context.Context, cache interface{}) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if _, ok := cache.(*Cache); !ok {
		cache = newCache(ctx)
	}

	return context.WithValue(ctx, ctxKey, cache)
}

func FromContext(ctx context.Context) *Cache {
	if cache, ok := ctx.Value(ctxKey).(*Cache); ok {
		return cache
	}

	return newCache(ctx)
}

func newCache(ctx context.Context) *Cache {
	return &Cache{
		logger: logger.FromContext(ctx),

		cacheMap: &syncmap.Map{},
	}
}

// Store stores value in cache
func (c *Cache) Store(key string, value interface{}) {
	c.logger.Info("store value in cache",
		zap.String("key", key),
		zap.Reflect("value", value),
	)

	c.cacheMap.Store(key, value)
}

// Load return value from cache
func (c *Cache) Load(key string) (interface{}, bool) {
	value, ok := c.cacheMap.Load(key)
	if !ok {
		return nil, false
	}

	c.logger.Info("load value from cache",
		zap.String("key", key),
		zap.Reflect("value", value),
	)

	return value, true
}
