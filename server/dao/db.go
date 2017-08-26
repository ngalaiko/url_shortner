package dao

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	// mandatory to init pq
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
	*sqlx.DB
	logger *logger.Logger
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
	cfg := config.FromContext(ctx)
	l := logger.FromContext(ctx)

	db := newDbHelper(cfg, l)

	if err := db.Healthcheck(ctx); err != nil {
		l.Error("error connection db",
			zap.Error(err),
		)
		time.Sleep(time.Second)
		return newDb(ctx)
	}

	l.Info("db connection created",
		zap.String("driver", cfg.Db.Driver),
		zap.String("config", cfg.Db.Connect),
	)

	return db
}

func newDbHelper(cfg *config.Config, l *logger.Logger) *Db {

	db, err := sqlx.Open(cfg.Db.Driver, cfg.Db.Connect)
	if err != nil {
		l.Panic("error while open db connection",
			zap.Error(err),
		)
	}

	db.SetMaxIdleConns(cfg.Db.MaxIdleConns)
	db.SetMaxOpenConns(cfg.Db.MaxOpenConns)

	l.Info("db connection created")

	return &Db{
		DB:     db,
		logger: l,
	}
}

// Healthcheck is a db healthcheck
func (db *Db) Healthcheck(ctx context.Context) error {
	_, err := db.Exec("SELECT 1")
	return err
}

// Mutate opens new tx, applies callback func and close tx
func (db *Db) Mutate(callback func(tx *Tx) error) error {
	tx := &Tx{
		db:     db.DB,
		logger: db.logger,
	}

	err := callback(tx)
	switch err {
	case nil:
		return tx.Commit()
	default:
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("error while rollback transaction: %s", err)
		}

		return err
	}
}
