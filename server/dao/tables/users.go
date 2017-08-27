// Code generated by generate_schema_tables.go DO NOT EDIT.

package tables

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/ngalayko/url_shortner/server/dao"
	"github.com/ngalayko/url_shortner/server/helpers"
	"github.com/ngalayko/url_shortner/server/schema"
)

// SelectUserById returns User from db or cache
func (t *Tables) SelectUserById(id uint64) (*schema.User, error) {
	ids := []uint64{id}

	uu, err := t.SelectUserByIds(ids)
	if err != nil {
		return nil, err
	}

	return uu[0], nil
}

// SelectUserByIds returns Users from db or cache
func (t *Tables) SelectUserByIds(ids []uint64) ([]*schema.User, error) {

	uu := make([]*schema.User, 0, len(ids))

	missingIds := make([]uint64, 0, len(ids))
	for _, id := range ids {
		value, ok := t.cache.Load(t.usersCacheKey(id))
		if !ok {
			missingIds = append(missingIds, id)
			continue
		}

		uu = append(uu, value.(*schema.User))
	}

	uuMissing := make([]*schema.User, 0, len(missingIds))
	if err := t.db.Select(uu,
		"SELECT *"+
			"FROM users"+
			"WHERE id IN ("+helpers.Uint64sToString(missingIds)+")",
	); err != nil {
		return nil, err
	}

	for _, uMissing := range uuMissing {
		uu = append(uu, uMissing)
		t.cache.Store(t.usersCacheKey(uMissing.ID), uMissing)
	}

	return uu, nil
}

// InsertUser inserts User in db and cache
func (t *Tables) InsertUser(u *schema.User) error {
	return t.db.Mutate(func(tx *dao.Tx) error {

		insertSQL := "INSERT INTO users" +
			"(first_name, last_name, created_at, deleted_at)" +
			"VALUES" +
			fmt.Sprintf("(%v, %v, %v, %v)",
				u.FirstName,
				u.LastName,
				u.CreatedAt,
				u.DeletedAt)

		_, err := tx.Exec(insertSQL)
		if err != nil {
			return err
		}

		t.logger.Info("User created",
			zap.Reflect("$.Name", u),
		)
		t.cache.Store(t.usersCacheKey(u.ID), u)
		return nil
	})
}

// UpdateUser updates User in db and cache
func (t *Tables) UpdateUser(u *schema.User) error {
	return t.db.Mutate(func(tx *dao.Tx) error {

		updateSQL := "UPDATE users" +
			"SET" +
			fmt.Sprintf("first_name = %v,", u.FirstName) +
			fmt.Sprintf("last_name = %v,", u.LastName) +
			fmt.Sprintf("created_at = %v,", u.CreatedAt) +
			fmt.Sprintf("deleted_at = %v", u.DeletedAt)

		_, err := tx.Exec(updateSQL)
		if err != nil {
			return err
		}

		t.logger.Info("User updated",
			zap.Reflect("$.Name", u),
		)
		t.cache.Store(t.usersCacheKey(u.ID), u)
		return nil
	})
}

func (t *Tables) usersCacheKey(id uint64) string {
	return fmt.Sprintf("User:%d", id)
}