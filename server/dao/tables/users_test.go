package tables

import (
	"fmt"
	"time"

	. "gopkg.in/check.v1"

	"github.com/ngalayko/url_shortner/server/dao"
	"github.com/ngalayko/url_shortner/server/schema"
)

func (s *TestTablesSuite) Test_InsertUser__should_insert_user(c *C) {
	user := s.testUser(c)

	c.Assert(user.ID, Not(Equals), uint64(0))
}

func (s *TestTablesSuite) Test_SelectUser__should_select_user(c *C) {
	user := s.testUser(c)

	selected, err := s.service.GetUserById(user.ID)
	if err != nil {
		c.Fatal(err)
	}

	c.Assert(user.ID, Equals, selected.ID)
	c.Assert(user.FirstName, Equals, selected.FirstName)
	c.Assert(user.LastName, Equals, selected.LastName)
}

func (s *TestTablesSuite) Test_SelectUsersByFields__should_select_links_by_params(c *C) {
	user1 := s.testUser(c)
	user2 := s.testUser(c)
	s.testLink(c)

	param1 := dao.NewParam(1).Add("id", user1.ID)
	param2 := dao.NewParam(1).Add("first_name", user2.FirstName)
	params := dao.NewParams(2).Append(param1).Append(param2)

	selected, err := s.service.SelectUsersByFields(params)
	if err != nil {
		c.Fatal(err)
	}

	c.Assert(len(selected), Equals, 2)
	c.Assert(selected[0].ID, Equals, user1.ID)
	c.Assert(selected[1].ID, Equals, user2.ID)
}

func (s *TestTablesSuite) Test_SelectUsersByFields__should_select_link_by_many_params(c *C) {
	user1 := s.testUser(c)

	user2 := s.testUser(c)
	user2.FirstName = user1.FirstName
	if err := s.service.UpdateUser(user2); err != nil {
		c.Fatal(err)
	}

	s.testUser(c)

	param1 := dao.NewParam(1).
		Add("id", user1.ID).
		Add("first_name", user1.FirstName)
	params := dao.NewParams(2).Append(param1)

	selected, err := s.service.SelectUsersByFields(params)
	if err != nil {
		c.Fatal(err)
	}

	c.Assert(len(selected), Equals, 1)
	c.Assert(selected[0].ID, Equals, user1.ID)
}

func (s *TestTablesSuite) Test_UpdateUser__should_update_user(c *C) {
	user := s.testUser(c)
	user.LastName = "updated"

	if err := s.service.UpdateUser(user); err != nil {
		c.Fatal(err)
	}

	updated, err := s.service.GetUserById(user.ID)
	if err != nil {
		c.Fatal(err)
	}

	c.Assert(user.ID, Equals, updated.ID)
	c.Assert(user.FirstName, Equals, updated.FirstName)
	c.Assert(user.LastName, Equals, updated.LastName)
}

func (s *TestTablesSuite) testUser(c *C) *schema.User {
	s.usersCount++

	user := &schema.User{
		FirstName:  fmt.Sprintf("first name %d", s.usersCount),
		LastName:   fmt.Sprintf("last name %d", s.usersCount),
		CreatedAt:  time.Now(),
		FacebookID: fmt.Sprintf("facebookID%d", s.usersCount),
	}

	if err := s.service.InsertUser(user); err != nil {
		c.Fatal(err)
	}

	return user
}
