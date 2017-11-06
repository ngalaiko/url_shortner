package users

import (
	"context"
	"database/sql"
	"time"

	"github.com/ngalayko/url_shortner/server/dao"
	"github.com/ngalayko/url_shortner/server/dao/tables"
	"github.com/ngalayko/url_shortner/server/facebook"
	"github.com/ngalayko/url_shortner/server/logger"
	"github.com/ngalayko/url_shortner/server/schema"
)

// Service is a users service
type Service struct {
	logger logger.ILogger
	tables *tables.Service

	facebookAPI *facebook.Api
}

func newUsers(ctx context.Context) *Service {
	u := &Service{
		logger: logger.FromContext(ctx),
		tables: tables.FromContext(ctx),

		facebookAPI: facebook.FromContext(ctx),
	}

	return u
}

// QueryUserById returns user by id
func (u *Service) QueryUserById(id uint64) (*schema.User, error) {

	user, err := u.tables.GetUserById(id)
	if err != nil {
		return nil, err
	}

	if user.DeletedAt != nil {
		return nil, sql.ErrNoRows
	}

	return user, nil
}

// QueryUserByFacebookUser returns user by facebook id
func (u *Service) QueryUserByFacebookUser(facebookUser *facebook.User) (*schema.User, error) {

	user, err := u.tables.GetUserByFields(dao.NewParam(1).Add("facebook_id", facebookUser.ID))
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

	if err := u.tables.InsertUser(user); err != nil {
		return nil, err
	}

	return user, nil
}
