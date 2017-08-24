package dao

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/ngalayko/url_shortner/server/config"
	"github.com/ngalayko/url_shortner/server/logger"
)

const (
	ctxKey dbCtxKey = "db_ctx_key"
)

type dbCtxKey string

type Db struct {
	*sqlx.DB
	logger *logger.Logger
}

func NewContext(ctx context.Context, db interface{}) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if _, ok := db.(*Db); !ok {
		db = newDb(ctx)
	}

	return context.WithValue(ctx, ctxKey, db)
}

func FromContext(ctx context.Context) *Db {
	if db, ok := ctx.Value(ctxKey).(*Db); ok {
		return db
	}

	return newDb(ctx)
}

func newDb(ctx context.Context) *Db {
	cfg := config.FromContext(ctx)
	l := logger.FromContext(ctx)

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

func (db *Db) Mutate(callback func(tx *sqlx.Tx) error) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	if err := callback(tx); err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("error while rollback transaction: %s", err)
		}

		return err
	}

	return tx.Commit()
}
