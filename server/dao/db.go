package dao

import (
	"context"
	"database/sql"
	"time"

	"go.uber.org/zap"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"

	_ "github.com/lib/pq"

	"github.com/ngalayko/url_shortner/server/config"
	"github.com/ngalayko/url_shortner/server/logger"
)

const (
	ctxKey dbCtxKey = "db_ctx_key"
)

type dbCtxKey string

// Db is a database service
type Db struct {
	*reform.DB
	config config.DbConfig
}

// NewContext stores Db in context
func NewContext(ctx context.Context, db interface{}) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if _, ok := db.(*Db); !ok {
		db = newDb(ctx)
	}

	return context.WithValue(ctx, ctxKey, db)
}

// FromContext return Db from context
func FromContext(ctx context.Context) *Db {
	if db, ok := ctx.Value(ctxKey).(*Db); ok {
		return db
	}

	return newDb(ctx)
}

func newDb(ctx context.Context) *Db {
	cfg := config.FromContext(ctx).Db
	l := logger.FromContext(ctx)

	db := newDbHelper(cfg, l)

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
		DB: reform.NewDB(conn, postgresql.Dialect, newDbLogger(logger)),
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
