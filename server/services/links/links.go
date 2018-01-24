package links

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/ngalayko/url_shortner/server/dao"
	"github.com/ngalayko/url_shortner/server/helpers"
	"github.com/ngalayko/url_shortner/server/logger"
	"github.com/ngalayko/url_shortner/server/schema"
)

const (
	defaultShortUrlLen = 6
	httpScheme         = "http"
)

// Service is a links service
type Service struct {
	ctx    context.Context
	logger logger.ILogger
	db     *dao.Db

	// channel of link ids to increment views of
	viewsQueue chan uint64

	// channel of links with updated user_id
	transferQueue chan *schema.Link
}

func newLinks(ctx context.Context) *Service {
	l := &Service{
		ctx:    ctx,
		logger: logger.FromContext(ctx),
		db:     dao.FromContext(ctx),

		viewsQueue:    make(chan uint64),
		transferQueue: make(chan *schema.Link),
	}

	go l.loop()

	return l
}

func (l *Service) loop() {
	for {
		select {
		case id := <-l.viewsQueue:
			if err := l.incrementNextLink(id); err != nil {
				l.logger.Error("error incrementing link",
					zap.Error(err),
				)
			}

		case link := <-l.transferQueue:
			if err := l.transferLink(link); err != nil {
				l.logger.Error("error while transfering link",
					zap.Error(err),
					zap.Reflect("link", link),
				)
			}

		case <-l.ctx.Done():
			return

		}
	}
}

func (l *Service) transferLink(link *schema.Link) error {
	if err := l.db.Update(link); err != nil {
		return err
	}

	l.logger.Info("link transfered",
		zap.Uint64("link id", link.ID),
		zap.Uint64("new user id", link.UserID),
	)
	return nil
}

func (l *Service) incrementNextLink(linkId uint64) error {
	link := &schema.Link{}
	if err := l.db.FindByPrimaryKeyTo(link, linkId); err != nil {
		return err
	}

	link.Views += 1
	if err := l.db.Update(link); err != nil {
		return err
	}

	l.logger.Info("link views incremented",
		zap.Uint64("id", link.ID),
		zap.Uint64("views", link.Views),
	)

	return nil
}

// TransferLinks transfer links from anon user to new user
func (l *Service) TransferLinks(userID uint64, links ...*schema.Link) error {
	if len(links) == 0 {
		return nil
	}

	for _, link := range links {
		link.UserID = userID
		l.transferQueue <- link
	}

	return nil
}

// CreateLink creates given link
func (l *Service) CreateLink(link *schema.Link) error {
	if err := prepareLink(link); err != nil {
		return err
	}

	existed, err := l.queryLinkByURlAndUserID(link.URL, link.UserID)
	switch {
	case err == sql.ErrNoRows:

	case err == nil:
		*link = *existed
		return nil

	case err != nil:
		return err

	}

	l.logger.Info("link created",
		zap.Reflect("link", link),
	)
	return l.db.Insert(link)
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

	return nil
}

// QueryLinkByShortUrl returns link by short url
func (l *Service) QueryLinkByShortUrl(shortUrl string) (*schema.Link, error) {
	link := &schema.Link{}
	if err := l.db.FindOneTo(link, "short_url", shortUrl); err != nil {
		return nil, err
	}

	if !link.Valid() {
		return nil, sql.ErrNoRows
	}

	l.viewsQueue <- link.ID

	l.logger.Info("link visited",
		zap.Reflect("link", link),
	)
	return link, nil
}

// QueryLinkByURlAndUserID returns link by userID and url
func (l *Service) queryLinkByURlAndUserID(url string, userID uint64) (*schema.Link, error) {
	link := &schema.Link{}
	tail := fmt.Sprintf(`
		WHERE
			url = %s AND
			user_id = %s AND
			deleted_at IS NULL
	`, l.db.Placeholder(1), l.db.Placeholder(2))
	if err := l.db.SelectOneTo(link, tail, url, userID); err != nil {
		return nil, err
	}

	if !link.Valid() {
		return nil, sql.ErrNoRows
	}

	if link.Anonim() {
		return nil, sql.ErrNoRows
	}

	return link, nil
}

// QueryLinkByShortUrl returns link by short url
func (l *Service) QueryLinksByUser(userID uint64) ([]*schema.Link, error) {
	rows, err := l.db.FindRows(schema.LinkTable, "user_id", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// todo: make big slice
	links := []*schema.Link{}
	for {
		link := &schema.Link{}
		if err := l.db.NextRow(link, rows); err != nil {
			break
		}
		if !link.Valid() {
			continue
		}

		links = append(links, link)
	}

	return links, nil
}

func (l *Service) deleteLink(link *schema.Link) error {
	if link.DeletedAt != nil {
		return nil
	}

	link.DeletedAt = new(time.Time)
	*link.DeletedAt = time.Now()

	l.logger.Info("link deleted",
		zap.Reflect("link", link),
	)
	return l.db.Update(link)
}
