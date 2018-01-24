package token

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	. "gopkg.in/check.v1"

	"github.com/ngalayko/url_shortner/server/config"
	"github.com/ngalayko/url_shortner/server/helpers"
	"github.com/ngalayko/url_shortner/server/schema"
)

type TestTokensSuite struct {
	ctx context.Context

	service *Service

	usersCount uint64
}

func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&TestTokensSuite{})

var suite *TestTokensSuite

func (s *TestTokensSuite) SetUpSuite(c *C) {
	suite = &TestTokensSuite{
		ctx: context.Background(),
	}
	s.ctx = config.NewContext(s.ctx, config.NewTestConfig())
	s.service = FromContext(s.ctx)
}

func (s *TestTokensSuite) Test_CreateUserToken__should_create_user_token(c *C) {
	user, err := s.createUser()
	c.Assert(err, IsNil)

	token, err := s.service.CreateUserToken(user)
	if err != nil {
		c.Fatal(err)
	}

	selected, err := s.service.GetUserToken(token.Token)
	if err != nil {
		c.Fatal(err)
	}

	c.Assert(token.UserID, Equals, user.ID)
	c.Assert(token.Token, Equals, selected.Token)
	c.Assert(token.ID, Equals, selected.ID)
}

func (s *TestTokensSuite) Test_DeleteUserToken__should_delete_user_token(c *C) {
	token, err := s.createToken()
	c.Assert(err, IsNil)

	if err := s.service.DeleteUserToken(token.UserID, token.Token); err != nil {
		c.Fatal(err)
	}

	_, err = s.service.GetUserToken(token.Token)
	c.Assert(err, Equals, sql.ErrNoRows)
}

func (s *TestTokensSuite) Test_GetUserToken__should_return_not_found_by_empty_token(c *C) {
	_, err := s.service.GetUserToken("")
	c.Assert(err, Equals, sql.ErrNoRows)
}

// helpers

func (s *TestTokensSuite) createToken() (*schema.UserToken, error) {
	user, err := s.createUser()
	if err != nil {
		return nil, err
	}

	return s.service.CreateUserToken(user)
}

func (s *TestTokensSuite) createUser() (*schema.User, error) {
	user := &schema.User{
		FirstName:  fmt.Sprintf("name %d", s.usersCount),
		LastName:   fmt.Sprintf("last name %d", s.usersCount),
		FacebookID: fmt.Sprintf("facebook id %s", helpers.RandomString(5)),
	}

	if err := s.service.db.Insert(user); err != nil {
		return nil, err
	}

	s.usersCount++
	return user, nil
}
