package db

import (
	"context"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"go.uber.org/zap"

	"github.com/ngalayko/url_shortner/server/config"
	"github.com/ngalayko/url_shortner/server/logger"
)

const (
	ctxKey dbCtxKey = "db_ctx_key"
)

type dbCtxKey string

type Db struct {
	Db *gorm.DB

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

	db, err := gorm.Open("postgres", cfg.Db.Connect)
	if err != nil {
		l.Fatal("cannot connect to postgres",
			zap.Error(err),
		)
	}

	l.Info("db connection created")

	return &Db{
		Db:     db,
		logger: l,
	}
}
