package tables

import (
	"context"

	"github.com/ngalayko/url_shortner/server/cache"
	"github.com/ngalayko/url_shortner/server/dao"
	"github.com/ngalayko/url_shortner/server/logger"
)

const (
	ctxKey tablesCtxKey = "tables_ctx_key"
)

type tablesCtxKey string

// Tables represents tables from db
type Tables struct {
	db     *dao.Db
	cache  cache.ICache
	logger logger.ILogger
}

// NewContext stores Tables in context
func NewContext(ctx context.Context, table interface{}) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if _, ok := table.(*Tables); !ok {
		table = newTables(ctx)
	}

	return context.WithValue(ctx, ctxKey, table)
}

// FromContext returns Tables from context
func FromContext(ctx context.Context) *Tables {
	if table, ok := ctx.Value(ctxKey).(*Tables); ok {
		return table
	}

	return newTables(ctx)
}

func newTables(ctx context.Context) *Tables {
	return &Tables{
		db:     dao.FromContext(ctx),
		cache:  cache.FromContext(ctx),
		logger: logger.FromContext(ctx),
	}
}
