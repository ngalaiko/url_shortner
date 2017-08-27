package dao

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/ngalayko/url_shortner/server/helpers"
	"github.com/ngalayko/url_shortner/server/logger"
)

// Tx is a tx wrapper over sqlx.Tx
type Tx struct {
	db     *sqlx.DB
	logger *logger.Logger

	tx *sqlx.Tx
	id string
}

func (t *Tx) begin() (*Tx, error) {
	if t.tx != nil {
		return t, nil
	}

	id, err := helpers.RandomString(10)
	if err != nil {
		return nil, err
	}

	tx, err := t.db.Beginx()
	if err != nil {
		return nil, err
	}

	t.id = id
	t.tx = tx

	t.logger.Debug("begin tx",
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

	t.logger.Debug("exec sql query",
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

	t.logger.Debug("exec sql query",
		zap.String("tx id", t.id),
		zap.String("query", query),
		zap.Reflect("args", args),
	)

	return t.tx.Exec(query, args...)
}

// Commit commits tx
func (t *Tx) Commit() error {
	t.logger.Debug("commit tx",
		zap.String("id", t.id),
	)

	return t.tx.Commit()
}

// Rollback rollbacks tx
func (t *Tx) Rollback() error {
	t.logger.Debug("begin tx",
		zap.String("id", t.id),
	)

	return t.tx.Rollback()
}
