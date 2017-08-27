// Code generated by generate_schema_tables.go DO NOT EDIT.

package tables

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/ngalayko/url_shortner/server/dao"
	"github.com/ngalayko/url_shortner/server/helpers"
	"github.com/ngalayko/url_shortner/server/schema"
)

// SelectLinkById returns Link from db or cache
func (t *Tables) SelectLinkById(id uint64) (*schema.Link, error) {
	ids := []uint64{id}

	ll, err := t.SelectLinkByIds(ids)
	if err != nil {
		return nil, err
	}

	return ll[0], nil
}

// SelectLinkByIds returns Links from db or cache
func (t *Tables) SelectLinkByIds(ids []uint64) ([]*schema.Link, error) {

	ll := make([]*schema.Link, 0, len(ids))

	missingIds := make([]uint64, 0, len(ids))
	for _, id := range ids {
		value, ok := t.cache.Load(t.linksCacheKey(id))
		if !ok {
			missingIds = append(missingIds, id)
			continue
		}

		ll = append(ll, value.(*schema.Link))
	}

	if len(missingIds) == 0 {
		return ll, nil
	}

	llMissing := make([]*schema.Link, 0, len(missingIds))
	if err := t.db.Select(&llMissing,
		"SELECT * "+
			"FROM links "+
			"WHERE id IN ("+helpers.Uint64sToString(missingIds)+")",
	); err != nil {
		return nil, err
	}

	for _, lMissing := range llMissing {
		ll = append(ll, lMissing)
		t.cache.Store(t.linksCacheKey(lMissing.ID), lMissing)
	}

	return ll, nil
}

// InsertLink inserts Link in db and cache
func (t *Tables) InsertLink(l *schema.Link) error {
	return t.db.Mutate(func(tx *dao.Tx) error {

		insertSQL := "INSERT INTO links " +
			"(user_id, url, short_url, clicks, views, expired_at, created_at, deleted_at) " +
			"VALUES " +
			"($1, $2, $3, $4, $5, $6, $7, $8) " +
			"RETURNING id"

		var id uint64
		if err := tx.Get(&id, insertSQL, l.UserID, l.URL, l.ShortURL, l.Clicks, l.Views, l.ExpiredAt, l.CreatedAt, l.DeletedAt); err != nil {
			return err
		}
		l.ID = id

		t.logger.Info("Link created",
			zap.Reflect("$.Name", l),
		)
		t.cache.Store(t.linksCacheKey(l.ID), l)
		return nil
	})
}

// UpdateLink updates Link in db and cache
func (t *Tables) UpdateLink(l *schema.Link) error {
	return t.db.Mutate(func(tx *dao.Tx) error {

		updateSQL := "UPDATE links " +
			"SET " +
			"user_id = $1, " +
			"url = $2, " +
			"short_url = $3, " +
			"clicks = $4, " +
			"views = $5, " +
			"expired_at = $6, " +
			"created_at = $7, " +
			"deleted_at = $8 " +
			fmt.Sprintf("WHERE id = %d", l.ID)

		_, err := tx.Exec(updateSQL, l.UserID, l.URL, l.ShortURL, l.Clicks, l.Views, l.ExpiredAt, l.CreatedAt, l.DeletedAt)
		if err != nil {
			return err
		}

		t.logger.Info("Link updated",
			zap.Reflect("$.Name", l),
		)
		t.cache.Store(t.linksCacheKey(l.ID), l)
		return nil
	})
}

func (t *Tables) linksCacheKey(id uint64) string {
	return fmt.Sprintf("Link:%d", id)
}
