package links

import (
	"context"
	"database/sql"
	"net/url"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/ngalayko/url_shortner/server/dao"
	"github.com/ngalayko/url_shortner/server/dao/tables"
	"github.com/ngalayko/url_shortner/server/helpers"
	"github.com/ngalayko/url_shortner/server/logger"
	"github.com/ngalayko/url_shortner/server/schema"
)

const (
	defaultExpire      = 24 * time.Hour
	defaultShortUrlLen = 6
	httpScheme         = "http"
)

// Links is a links service
type Links struct {
	logger logger.ILogger
	tables *tables.Tables

	//channel of link ids to increment views of
	viewsQueue chan uint64
}

func newLinks(ctx context.Context) *Links {
	l := &Links{
		logger: logger.FromContext(ctx),
		tables: tables.FromContext(ctx),

		viewsQueue: make(chan uint64),
	}

	go l.loop()

	return l
}

func (l *Links) loop() {
	for id := range l.viewsQueue {
		if err := l.incrementNextLink(id); err != nil {
			l.logger.Error("error incrementing link",
				zap.Error(err),
			)
		}
	}
}

func (l *Links) incrementNextLink(linkId uint64) error {
	link, err := l.tables.GetLinkById(linkId)
	if err != nil {
		return err
	}

	link.Views += 1
	if err := l.tables.UpdateLink(link); err != nil {
		return err
	}

	l.logger.Info("link views incremented",
		zap.Uint64("id", link.ID),
		zap.Uint64("views", link.Views),
	)

	return nil
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

	link, err := l.tables.GetLinkByFields(dao.NewParam(1).Add("short_url", shortUrl))
	if err != nil {
		return nil, err
	}

	if !link.Valid() {
		return nil, sql.ErrNoRows
	}

	l.viewsQueue <- link.ID

	return link, nil
}

// QueryLinkByShortUrl returns link by short url
func (l *Links) QueryLinksByUser(userID uint64) ([]*schema.Link, error) {
	return l.tables.SelectLinksByFields(dao.NewParams(1).Append(dao.NewParam(1).Add("user_id", userID)))
}
