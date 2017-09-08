package links

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/ngalayko/url_shortner/server/dao/tables"
	"github.com/ngalayko/url_shortner/server/helpers"
	"github.com/ngalayko/url_shortner/server/logger"
	"github.com/ngalayko/url_shortner/server/schema"
	"strings"
)

const (
	defaultExpire      = 24 * time.Hour
	defaultShortUrlLen = 6
	httpScheme         = "http"
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

	if err := prepareLink(link); err != nil {
		return err
	}

	return l.tables.InsertLink(link)
}

func prepareLink(link *schema.Link) error {
	if !strings.HasPrefix(link.URL, httpScheme) {
		link.URL = httpScheme + "://" + link.URL
	}

	uri, err := url.ParseRequestURI(link.URL)
	if err != nil {
		return err
	}

	now := time.Now()
	link.URL = uri.String()
	link.CreatedAt = now
	link.ShortURL = helpers.RandomString(defaultShortUrlLen)
	link.ExpiredAt = now.Add(defaultExpire)

	return nil
}

// QueryLinkByShortUrl returns link by short url
func (l *Links) QueryLinkByShortUrl(shortUrl string) (*schema.Link, error) {

	link, err := l.tables.SelectLinkByFields(map[string]interface{}{"short_url": shortUrl})
	if err != nil {
		return nil, err
	}

	if link.ExpiredAt.Before(time.Now()) {
		return link, fmt.Errorf("Link has expired at %s", link.ExpiredAt)
	}

	return link, nil
}
