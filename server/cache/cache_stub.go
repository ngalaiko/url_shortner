package cache

import "fmt"

type CacheStub struct{}

// NewStubCache returns stub cache
func NewStubCache() *CacheStub {
	fmt.Printf("stub cache")
	return &CacheStub{}
}

// Store does nothing
func (cs *CacheStub) Store(key string, value interface{}) {}

// Load returns nothing
func (cs *CacheStub) Load(key string) (interface{}, bool) { return nil, false }
