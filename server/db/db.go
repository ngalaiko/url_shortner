package db

import (
	"context"
	"database/sql"
	"time"

	"go.uber.org/zap"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"

	// connector for postgreSQL
	_ "github.com/lib/pq"

	"github.com/ngalayko/url_shortner/server/cache"
	"github.com/ngalayko/url_shortner/server/config"
	"github.com/ngalayko/url_shortner/server/logger"
)

const (
	ctxKey dbCtxKey = "db_ctx_key"
)

type dbCtxKey string

// Db is a database service
type Db struct {
	db    *reform.DB
	cache cache.ICache
}

func newDb(ctx context.Context) *Db {
	cfg := config.FromContext(ctx).Db
	l := logger.FromContext(ctx)

	db := newDbHelper(cfg, l)
	db.cache = cache.FromContext(ctx)

	l.Info("db connection created",
		zap.String("driver", cfg.Driver),
		zap.String("config", cfg.Connect),
	)

	return db
}

func newDbHelper(cfg config.DbConfig, logger logger.ILogger) *Db {
	conn, err := sql.Open(cfg.Driver, cfg.Connect)
	if err != nil {
		logger.Panic("error while open db connection",
			zap.Error(err),
		)
	}
	logger.Info("db connection created")

	conn.SetMaxIdleConns(cfg.MaxIdleConns)
	conn.SetMaxOpenConns(cfg.MaxOpenConns)

	return &Db{
		db: reform.NewDB(conn, postgresql.Dialect, newDbLogger(logger)),
	}
}

type dbLogger struct {
	logger logger.ILogger
}

func newDbLogger(logger logger.ILogger) *dbLogger {
	return &dbLogger{
		logger: logger,
	}
}

func (l *dbLogger) Before(query string, args []interface{}) {
	l.logger.Debug(
		"execute query",
		zap.String("query", query),
		zap.Reflect("args", args),
	)
}

func (l *dbLogger) After(query string, args []interface{}, d time.Duration, err error) {
	l.logger.Debug(
		"query finished",
		zap.String("query", query),
		zap.Reflect("args", args),
		zap.Duration("duration", d),
		zap.Error(err),
	)
}
