package cache

// stub us stub cache
type stub struct{}

// NewStubCache returns stub cache
func NewStubCache() *stub {
	return &stub{}
}

// Store does nothing
func (cs *stub) Store(key string, value interface{}) {}

// Load returns nothing
func (cs *stub) Load(key string) (interface{}, bool) { return nil, false }
