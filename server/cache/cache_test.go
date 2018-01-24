package cache

import (
	"context"
	"testing"

	. "gopkg.in/check.v1"

	"github.com/ngalayko/url_shortner/server/logger"
)

func Benchmark_Store(b *testing.B) {

	ctx := logger.NewContext(nil, nil)
	cache := FromContext(ctx)

	for i := 0; i < b.N; i++ {
		cache.Store("key", "val")
	}
}

type TestCacheSuite struct {
	cache ICache
}

func Test(t *testing.T) { TestingT(t) }

var suite *TestCacheSuite

var _ = Suite(&TestCacheSuite{})

func (s *TestCacheSuite) SetUpSuite(c *C) {
	s.cache = FromContext(context.Background())
}

func (s *TestCacheSuite) Test_Load(c *C) {
	key, value := "key", "value"

	s.cache.Store(key, value)

	v, ok := s.cache.Load(key)
	c.Assert(ok, Equals, true)
	c.Assert(v, Equals, value)
}
