package user_token

import (
	"context"
	"database/sql"
	"time"

	"github.com/ngalayko/url_shortner/server/dao"
	"github.com/ngalayko/url_shortner/server/helpers"
	"github.com/ngalayko/url_shortner/server/logger"
	"github.com/ngalayko/url_shortner/server/schema"
)

const (
	defaultTokenLength = 15
	defaultExpiredTime = time.Hour * 24 * 30 // month
)

// Service is a user token service
type Service struct {
	logger logger.ILogger
	db     *dao.Db
}

func newTokens(ctx context.Context) *Service {
	return &Service{
		logger: logger.FromContext(ctx),
		db:     dao.FromContext(ctx),
	}
}

// GetUserToken returns token
func (t *Service) GetUserToken(token string) (*schema.UserToken, error) {

	if len(token) == 0 {
		return nil, sql.ErrNoRows
	}

	userToken := &schema.UserToken{}
	err := t.db.FindOneTo(userToken, "token", token)
	if err != nil {
		return nil, err
	}

	if userToken.ExpiredAt.Before(time.Now()) {
		return nil, sql.ErrNoRows
	}

	return userToken, nil
}

// CreateUserToken created userToken for user
func (t *Service) CreateUserToken(user *schema.User) (*schema.UserToken, error) {

	token := &schema.UserToken{
		UserID:    user.ID,
		ExpiredAt: time.Now().Add(defaultExpiredTime),
		Token:     helpers.RandomString(defaultTokenLength),
	}

	if err := t.db.Insert(token); err != nil {
		return nil, err
	}

	return token, nil
}