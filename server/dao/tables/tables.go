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

// Service represents tables from db
type Service struct {
	db     *dao.Db
	cache  cache.ICache
	logger logger.ILogger
}

// NewContext stores Service in context
func NewContext(ctx context.Context, table interface{}) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if _, ok := table.(*Service); !ok {
		table = newTables(ctx)
	}

	return context.WithValue(ctx, ctxKey, table)
}

// FromContext returns Service from context
func FromContext(ctx context.Context) *Service {
	if table, ok := ctx.Value(ctxKey).(*Service); ok {
		return table
	}

	return newTables(ctx)
}

func newTables(ctx context.Context) *Service {
	return &Service{
		db:     dao.FromContext(ctx),
		cache:  cache.FromContext(ctx),
		logger: logger.FromContext(ctx),
	}
}
