package users

import (
	"context"
	"fmt"
	"testing"

	. "gopkg.in/check.v1"

	"github.com/ngalayko/url_shortner/server/config"
	"github.com/ngalayko/url_shortner/server/helpers"
	"github.com/ngalayko/url_shortner/server/schema"
)

type TestUsersSuite struct {
	ctx context.Context

	service *Service

	usersCount uint64
}

func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&TestUsersSuite{})

var suite *TestUsersSuite

func (s *TestUsersSuite) SetUpSuite(c *C) {
	suite = &TestUsersSuite{
		ctx: context.Background(),
	}
	s.ctx = config.NewContext(s.ctx, config.NewTestConfig())
	s.service = FromContext(s.ctx)
}

func (s *TestUsersSuite) Test_QueryUserById__should_return_user_by_id(c *C) {
	user, err := s.createUser()
	c.Assert(err, IsNil)

	selected, err := s.service.QueryUserByID(user.ID)
	if err != nil {
		c.Fatal(err)
	}

	c.Assert(selected.ID, Equals, user.ID)
}

// helpers

func (s *TestUsersSuite) createUser() (*schema.User, error) {
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
