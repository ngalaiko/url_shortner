package cache

type Stub struct{}

// NewStubCache returns stub cache
func NewStubCache() *Stub {
	return &Stub{}
}

// Store does nothing
func (cs *Stub) Store(key string, value interface{}) {}

// Load returns nothing
func (cs *Stub) Load(key string) (interface{}, bool) { return nil, false }
