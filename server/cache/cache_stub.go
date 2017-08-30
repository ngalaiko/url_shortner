package cache

import "fmt"

type cacheStub struct{}

func NewStubCache() *cacheStub {
	fmt.Printf("stub cache")
	return &cacheStub{}
}

// Store does nothing
func (cs *cacheStub) Store(key string, value interface{}) {}

// Load returns nothing
func (cs *cacheStub) Load(key string) (interface{}, bool) { return nil, false }
