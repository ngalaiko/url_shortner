package cache

import (
	"context"
	"time"

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
	ctx    context.Context
	logger logger.ILogger

	operationsChan chan operation
}

type operation func(map[string]interface{})

func newCache(ctx context.Context) *Cache {
	c := &Cache{
		ctx:            ctx,
		logger:         logger.FromContext(ctx),
		operationsChan: make(chan operation),
	}
	go c.background()

	return c
}

// Store stores value in cache
func (c *Cache) Store(key string, value interface{}) {
	start := time.Now()

	c.store(key, value)

	c.logger.Debug("store value in cache",
		zap.String("key", key),
		zap.Reflect("value", value),
		zap.Duration("duration", time.Since(start)),
	)
}

// Load return value from cache
func (c *Cache) Load(key string) (interface{}, bool) {
	start := time.Now()

	value, ok := c.load(key)
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

func (c *Cache) load(key string) (interface{}, bool) {
	resultChan := make(chan interface{})
	c.operationsChan <- func(cacheMap map[string]interface{}) {
		resultChan <- cacheMap[key]
	}

	result := <-resultChan
	return result, !(result == nil)
}

func (c *Cache) store(key string, value interface{}) {
	c.operationsChan <- func(cacheMap map[string]interface{}) {
		cacheMap[key] = value
	}
}

func (c *Cache) background() {
	cacheMap := map[string]interface{}{}

	for {
		select {
		case operation := <-c.operationsChan:
			operation(cacheMap)

		case <-c.ctx.Done():
			return
		}
	}
}
