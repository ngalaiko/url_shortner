package session

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/ngalayko/url_shortner/server/cache"
	"github.com/ngalayko/url_shortner/server/helpers"
	"github.com/ngalayko/url_shortner/server/logger"
	"github.com/ngalayko/url_shortner/server/schema"
)

var (
	ErrorNoSuchSession = errors.New("no such session")
)

type ISession interface {
	Create() *schema.Session
	Load(key string) (*schema.Session, error)
	Update(session *schema.Session) error
}

type service struct {
	logger logger.ILogger
	cache  cache.ICache
}

func newSession(ctx context.Context) *service {
	return &service{
		logger: logger.FromContext(ctx),
		cache:  cache.FromContext(ctx),
	}
}

// Create creates new session
func (s *service) Create() *schema.Session {
	session := &schema.Session{
		Key:     helpers.RandomString(10),
		LinkIDs: []uint64{},
	}

	s.cache.Store(session.Key, session)
	s.logger.Info("session created",
		zap.String("key", session.Key),
	)
	return session
}

// Update updates session
func (s *service) Update(session *schema.Session) error {
	if session == nil || session.Key == "" {
		return errors.New("can't update empty session")
	}

	s.cache.Store(session.Key, session)
	s.logger.Info("session updated",
		zap.String("key", session.Key),
		zap.Uint64s("link ids", session.LinkIDs),
	)

	return nil
}

// Load returns session by key
func (s *service) Load(key string) (*schema.Session, error) {
	session, ok := s.cache.Load(key)
	if !ok {
		return nil, ErrorNoSuchSession
	}

	if s, ok := session.(*schema.Session); ok {
		return s, nil
	}
	return nil, errors.New("error while casing session")
}
