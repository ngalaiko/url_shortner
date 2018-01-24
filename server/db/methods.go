package db

import (
	"database/sql"
	"fmt"
	"reflect"

	reform "gopkg.in/reform.v1"
)

// Insert inserts record
func (t *Db) Insert(record reform.Record) error {
	if err := t.db.Save(record); err != nil {
		return err
	}

	t.cache.Store(cacheKey(record), record)

	return nil
}

// Update updates record
func (t *Db) Update(record reform.Record) error {
	if err := t.db.Update(record); err != nil {
		return err
	}

	t.cache.Store(cacheKey(record), record)

	return nil
}

// FindByPrimaryKeyTo finds first recoed by pk
func (t *Db) FindByPrimaryKeyTo(record reform.Record, pk interface{}) error {
	if fromCache, ok := t.cache.Load(cacheKeyWithValue(record, pk)); ok {
		cacheValue := reflect.ValueOf(fromCache).Elem()
		recordValue := reflect.ValueOf(record).Elem()
		recordValue.Set(cacheValue)
		return nil
	}
	return t.db.FindByPrimaryKeyTo(record, pk)
}

// FindOneTo finds first record by column value
func (t *Db) FindOneTo(str reform.Struct, column string, arg interface{}) error {
	return t.db.FindOneTo(str, column, arg)
}

// FindAllFromByPk returns all records from db By field
func (t *Db) FindAllFromByPk(view reform.View, pkColumn string, pks ...interface{}) ([]reform.Struct, error) {
	fromCache := make([]reform.Struct, 0, 0)
	for _, pk := range pks {
		found, ok := t.cache.Load(cacheKeyStringInterface(pkColumn, pk))
		if !ok {
			continue
		}

		if foundStruct, ok := found.(reform.Struct); ok {
			fromCache = append(fromCache, foundStruct)
		}
	}
	if len(fromCache) == len(pks) {
		return fromCache, nil
	}

	fromDb, err := t.db.FindAllFrom(view, pkColumn, pks...)
	if err != nil {
		return nil, err
	}
	return append(fromCache, fromDb...), nil
}

// Exec executes query with args
func (t *Db) Exec(query string, args ...interface{}) (sql.Result, error) {
	return t.db.Exec(query, args...)
}

// Placeholder returns placeholder for certain sql dialect
func (t *Db) Placeholder(pos int) string {
	return t.db.Placeholder(pos)
}

// SelectOneTo selects record by custom query
func (t *Db) SelectOneTo(str reform.Struct, tail string, args ...interface{}) error {
	return t.db.SelectOneTo(str, tail, args...)
}

// FindRows return rows by column
func (t *Db) FindRows(view reform.View, column string, arg interface{}) (*sql.Rows, error) {
	return t.db.FindRows(view, column, arg)
}

// SelectRows return sql.Rows by query
func (t *Db) SelectRows(view reform.View, tail string, args ...interface{}) (*sql.Rows, error) {
	return t.db.SelectRows(view, tail, args...)
}

// NextRow return next structure from rows
func (t *Db) NextRow(str reform.Struct, rows *sql.Rows) error {
	return t.db.NextRow(str, rows)
}

// InTransaction wraps function with begin : commit ? rollback
func (t *Db) InTransaction(f func(*reform.TX) error) error {
	return t.db.InTransaction(f)
}

//
// helpers
//

func cacheKeyStringInterface(name string, val interface{}) string {
	return fmt.Sprintf("%s%v", name, val)
}

func cacheKeyWithValue(record reform.Record, pk interface{}) string {
	return cacheKeyStringInterface(record.View().Name(), pk)
}

func cacheKey(record reform.Record) string {
	return cacheKeyWithValue(record, record.PKValue())
}
