package users

import (
	"context"
	"database/sql"

	"github.com/ngalayko/url_shortner/server/dao/tables"
	"github.com/ngalayko/url_shortner/server/logger"
	"github.com/ngalayko/url_shortner/server/schema"
)

// Users is a users service
type Users struct {
	logger *logger.Logger
	tables *tables.Tables
}

func newUsers(ctx context.Context) *Users {
	u := &Users{
		logger: logger.FromContext(ctx).Prefix("users"),
		tables: tables.FromContext(ctx),
	}

	return u
}

// QueryUserById returns user by id
func (u *Users) QueryUserById(id uint64) (*schema.User, error) {

	user, err := u.tables.GetUserById(id)
	if err != nil {
		return nil, err
	}

	if user.DeletedAt != nil {
		return nil, sql.ErrNoRows
	}

	return user, nil
}
