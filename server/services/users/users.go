package users

import (
	"context"
	"database/sql"
	"time"

	"go.uber.org/zap"

	"github.com/ngalayko/url_shortner/server/db"
	"github.com/ngalayko/url_shortner/server/facebook"
	"github.com/ngalayko/url_shortner/server/logger"
	"github.com/ngalayko/url_shortner/server/schema"
)

// Service is a users service
type Service struct {
	logger logger.ILogger
	db     *db.Db

	facebookAPI *facebook.API
}

func newUsers(ctx context.Context) *Service {
	u := &Service{
		logger: logger.FromContext(ctx),
		db:     db.FromContext(ctx),

		facebookAPI: facebook.FromContext(ctx),
	}

	return u
}

// QueryUserByID returns user by id
func (u *Service) QueryUserByID(id uint64) (*schema.User, error) {
	user := &schema.User{}
	if err := u.db.FindByPrimaryKeyTo(user, id); err != nil {
		return nil, err
	}

	if user.DeletedAt != nil {
		return nil, sql.ErrNoRows
	}

	return user, nil
}

// QueryUserByFacebookUser returns user by facebook id
func (u *Service) QueryUserByFacebookUser(facebookUser *facebook.User) (*schema.User, error) {
	user := &schema.User{}
	err := u.db.FindOneTo(user, "facebook_id", facebookUser.ID)
	switch {
	case err == nil:
		return user, nil

	case err == sql.ErrNoRows:

	case err != nil:
		return nil, err

	}

	user = &schema.User{
		FirstName:  facebookUser.FirstName,
		LastName:   facebookUser.LastName,
		FacebookID: facebookUser.ID,
		CreatedAt:  time.Now(),
	}

	if err := u.db.Insert(user); err != nil {
		return nil, err
	}

	u.logger.Info("user created",
		zap.Reflect("user", user),
	)
	return user, nil
}
