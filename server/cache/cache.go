package cache

import (
	"context"
	"time"

	"golang.org/x/sync/syncmap"

	"github.com/ngalayko/url_shortner/server/logger"
	"go.uber.org/zap"
)

const (
	ctxKey cacheContextKey = "cache_ctx_key"
)

type cacheContextKey string

// ICache is a cache interface
type ICache interface {
	Store(key string, value interface{})
	Load(key string) (interface{}, bool)
}

// Cache is a cache service
type Cache struct {
	logger logger.ILogger

	cacheMap *syncmap.Map
}

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

func newCache(ctx context.Context) *Cache {
	return &Cache{
		logger:   logger.FromContext(ctx),
		cacheMap: &syncmap.Map{},
	}
}

// Store stores value in cache
func (c *Cache) Store(key string, value interface{}) {
	start := time.Now()

	if _, ok := c.cacheMap.Load(key); ok {
		return
	}

	c.cacheMap.Store(key, value)

	c.logger.Debug("store value in cache",
		zap.String("key", key),
		zap.Reflect("value", value),
		zap.Duration("duration", time.Since(start)),
	)
}

// Load return value from cache
func (c *Cache) Load(key string) (interface{}, bool) {
	start := time.Now()

	value, ok := c.cacheMap.Load(key)
	if !ok {
		return nil, false
	}

	c.logger.Debug("load value from cache",
		zap.String("key", key),
		zap.Reflect("value", value),
		zap.Duration("duration", time.Since(start)),
	)

	return value, true
}
