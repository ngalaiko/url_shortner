package cache

import (
	"testing"

	"github.com/ngalayko/url_shortner/server/logger"
)

func Benchmark_Store(b *testing.B) {

	ctx := logger.NewContext(nil, nil)
	cache := FromContext(ctx)

	for i := 0; i < b.N; i ++ {
		cache.Store("key", "val")
	}

}
