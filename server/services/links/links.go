package links

import (
	"context"
	"time"

	"github.com/ngalayko/url_shortner/server/dao/tables"
	"github.com/ngalayko/url_shortner/server/logger"
	"github.com/ngalayko/url_shortner/server/schema"
)

const (
	defaultExpire = 24 * time.Hour
)

// Links is a links service
type Links struct {
	logger *logger.Logger
	tables *tables.Tables
}

func newLinks(ctx context.Context) *Links {
	return &Links{
		logger: logger.FromContext(ctx),
		tables: tables.FromContext(ctx),
	}
}

// CreateLink creates given link
func (l *Links) CreateLink(link *schema.Link) error {

	now := time.Now()
	link.CreatedAt = now
	link.ShortURL = link.URL + "_short"
	link.ExpiredAt = now.Add(defaultExpire)

	if err := l.tables.InsertLink(link); err != nil {
		return err
	}

	return nil
}
