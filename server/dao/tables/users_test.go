package tables

import (
	"fmt"
	"time"

	. "gopkg.in/check.v1"

	"github.com/ngalayko/url_shortner/server/schema"
)

func (s *TestTablesSuite) Test_InsertUser__should_insert_user(c *C) {
	user := s.testUser(c)

	c.Assert(user.ID, Not(Equals), uint64(0))
}

func (s *TestTablesSuite) Test_SelectUser__should_select_user(c *C) {
	user := s.testUser(c)

	selected, err := s.service.SelectUserById(user.ID)
	if err != nil {
		c.Fatal(err)
	}

	c.Assert(user, DeepEquals, selected)
}

func (s *TestTablesSuite) Test_UpdateUser__should_update_user(c *C) {
	user := s.testUser(c)
	user.LastName = "updated"

	if err := s.service.UpdateUser(user); err != nil {
		c.Fatal(err)
	}

	updated, err := s.service.SelectUserById(user.ID)
	if err != nil {
		c.Fatal(err)
	}

	c.Assert(user, DeepEquals, updated)
}

func (s *TestTablesSuite) testUser(c *C) *schema.User {
	s.usersCount++

	user := &schema.User{
		FirstName: fmt.Sprintf("first name %d", s.usersCount),
		LastName:  fmt.Sprintf("last name %d", s.usersCount),
		CreatedAt: time.Now(),
	}

	if err := s.service.InsertUser(user); err != nil {
		c.Fatal(err)
	}

	return user
}
