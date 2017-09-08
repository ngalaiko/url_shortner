package dao

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/ngalayko/url_shortner/server/helpers"
	"github.com/ngalayko/url_shortner/server/logger"
	"time"
)

// Tx is a tx wrapper over sqlx.Tx
type Tx struct {
	db     *sqlx.DB
	logger *logger.Logger

	tx *sqlx.Tx
	id string

	start time.Time
}

func (t *Tx) begin() (*Tx, error) {
	if t.tx != nil {
		return t, nil
	}

	id := helpers.RandomString(10)

	tx, err := t.db.Beginx()
	if err != nil {
		return nil, err
	}

	t.id = id
	t.tx = tx
	t.start = time.Now()

	t.logger.Info("begin tx",
		zap.String("id", t.id),
	)

	return t, nil
}

// Get executes named sql
func (t *Tx) Get(dest interface{}, query string, args ...interface{}) error {
	t, err := t.begin()
	if err != nil {
		return err
	}

	t.logger.Info("exec sql query",
		zap.String("tx id", t.id),
		zap.String("query", query),
		zap.Reflect("args", args),
	)

	return t.tx.Get(dest, query, args...)
}

// Exec executes sql
func (t *Tx) Exec(query string, args ...interface{}) (sql.Result, error) {
	t, err := t.begin()
	if err != nil {
		return nil, err
	}

	start := time.Now()

	defer t.logger.Info("exec sql query",
		zap.String("tx id", t.id),
		zap.String("query", query),
		zap.Reflect("args", args),
		zap.Duration("duration", time.Since(start)),
	)

	return t.tx.Exec(query, args...)
}

// Commit commits tx
func (t *Tx) Commit() error {
	defer t.logger.Info("commit tx",
		zap.String("id", t.id),
		zap.Duration("duration", time.Since(t.start)),
	)

	return t.tx.Commit()
}

// Rollback rollbacks tx
func (t *Tx) Rollback() error {
	defer t.logger.Info("begin tx",
		zap.String("id", t.id),
		zap.Duration("duration", time.Since(t.start)),
	)

	return t.tx.Rollback()
}
