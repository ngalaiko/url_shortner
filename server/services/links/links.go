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

// Service is a links service
type Service struct {
	logger logger.ILogger
	tables *tables.Service

	//channel of link ids to increment views of
	viewsQueue chan uint64
}

func newLinks(ctx context.Context) *Service {
	l := &Service{
		logger: logger.FromContext(ctx),
		tables: tables.FromContext(ctx),

		viewsQueue: make(chan uint64),
	}

	go l.loop()

	return l
}

func (l *Service) loop() {
	for id := range l.viewsQueue {
		if err := l.incrementNextLink(id); err != nil {
			l.logger.Error("error incrementing link",
				zap.Error(err),
			)
		}
	}
}

func (l *Service) incrementNextLink(linkId uint64) error {
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
func (l *Service) CreateLink(link *schema.Link) (*schema.Link, error) {

	if err := prepareLink(link); err != nil {
		return nil, err
	}

	existed, err := l.QueryLinkByURlAndUserID(link.URL, link.UserID)
	switch {
	case err == sql.ErrNoRows:

	case err == nil:
		return existed, nil

	case err != nil:
		return nil, err

	}

	if err := l.tables.InsertLink(link); err != nil {
		return nil, err
	}

	return link, nil
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
func (l *Service) QueryLinkByShortUrl(shortUrl string) (*schema.Link, error) {

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

// QueryLinkByURlAndUserID returns link by uset id and url
func (l *Service) QueryLinkByURlAndUserID(url string, userID uint64) (*schema.Link, error) {
	link, err := l.tables.GetLinkByFields(dao.NewParam(2).Add("url", url).Add("user_id", userID))
	if err != nil {
		return nil, err
	}

	if !link.Valid() {
		return nil, sql.ErrNoRows
	}

	return link, nil
}

// QueryLinkByShortUrl returns link by short url
func (l *Service) QueryLinksByUser(userID uint64) ([]*schema.Link, error) {
	return l.tables.SelectLinksByFields(dao.NewParams(1).Append(dao.NewParam(1).Add("user_id", userID)))
}
