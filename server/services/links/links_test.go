package links

import (
	"context"
	"fmt"
	"testing"

	. "gopkg.in/check.v1"

	"github.com/ngalayko/url_shortner/server/config"
	"github.com/ngalayko/url_shortner/server/helpers"
	"github.com/ngalayko/url_shortner/server/schema"
)

type TestLinksSuite struct {
	ctx context.Context

	service *Service

	linksCount uint64
	usersCount uint64
}

func Test(t *testing.T) { TestingT(t) }

var suite *TestLinksSuite

var _ = Suite(&TestLinksSuite{})

func (s *TestLinksSuite) SetUpSuite(c *C) {
	suite = &TestLinksSuite{
		ctx: context.Background(),
	}
	s.ctx = config.NewContext(s.ctx, config.NewTestConfig())
	s.service = FromContext(s.ctx)
}

func (s *TestLinksSuite) Test_CreateLink__should_create_link(c *C) {
	_, err := s.createLink()
	c.Assert(err, IsNil)
}

func (s *TestLinksSuite) Test_CreateLink__should_return_new_link_for_anon_user(c *C) {
	link1, err := s.createLink(
		withUserID(0),
	)
	if err != nil {
		c.Fatal(err)
	}

	link2, err := s.createLink(
		withUserID(0),
		withUrl(link1.URL),
	)
	if err != nil {
		c.Fatal(err)
	}

	c.Assert(link1.ID, Not(Equals), link2.ID)
}

func (s *TestLinksSuite) Test_CreateLink__should_return_same_link_if_exists_for_not_anon_user(c *C) {
	userId := uint64(1)

	link1, err := s.createLink(
		withUserID(userId),
	)
	if err != nil {
		c.Fatal(err)
	}

	link2, err := s.createLink(
		withUrl(link1.URL),
		withUserID(userId),
	)
	if err != nil {
		c.Fatal(err)
	}

	c.Assert(link1.ID, Equals, link2.ID)
}

func (s *TestLinksSuite) Test_prepareLink__should_set_schema(c *C) {
	url := "vk.com"

	link := &schema.Link{
		URL: url,
	}

	if err := prepareLink(link); err != nil {
		c.Fatal(err)
	}

	c.Assert("http://"+url, Equals, link.URL)
}

func (s *TestLinksSuite) Test_prepareLink__should_not_set_schema(c *C) {
	url := "https://vk.com"

	link := &schema.Link{
		URL: url,
	}

	if err := prepareLink(link); err != nil {
		c.Fatal(err)
	}

	c.Assert(link.URL, Equals, url)
}

func (s *TestLinksSuite) Test_prepareLink__should_create_short_uri(c *C) {
	link := &schema.Link{
		URL: "vk.com",
	}

	if err := prepareLink(link); err != nil {
		c.Fatal(err)
	}

	c.Assert(link.ShortURL, Not(Equals), "")
}

func (s *TestLinksSuite) Test_prepareLink__should_not_validate_url(c *C) {
	link := &schema.Link{
		URL: "http//vk.com",
	}

	err := prepareLink(link)
	c.Assert(err, NotNil)
}

func (s *TestLinksSuite) Test_QueryLinksByUser__should_not_query_deleted_links(c *C) {
	link, err := s.createLink()
	c.Assert(err, IsNil)

	deletedLink, err := s.createLink(
		withUserID(link.UserID),
	)
	c.Assert(err, IsNil)

	if err := s.service.deleteLink(deletedLink); err != nil {
		c.Fatal(err)
	}

	links, err := s.service.QueryLinksByUser(link.UserID)
	c.Assert(err, IsNil)

	c.Assert(1, Equals, len(links))
}

func (s *TestLinksSuite) Test_TransferLinks__should_change_links_owner(c *C) {
	user, err := s.createUser()
	c.Assert(err, IsNil)

	link, err := s.createLink()
	c.Assert(err, IsNil)

	if err := s.service.TransferLinks(user.ID, link.ID); err != nil {
		c.Fatal(err)
	}

	links, err := s.service.QueryLinksByUser(user.ID)
	c.Assert(err, IsNil)

	for _, link := range links {
		c.Assert(link.UserID, Equals, user.ID)
	}
}

// helpers

func (s *TestLinksSuite) createLink(opts ...optionFunc) (*schema.Link, error) {
	user, err := s.createUser()
	if err != nil {
		return nil, err
	}

	link := &schema.Link{
		UserID: user.ID,
		URL:    fmt.Sprintf("http://vk.com/%d", s.linksCount),
	}

	for _, opt := range opts {
		opt(link)
	}

	if err := s.service.CreateLink(link); err != nil {
		return nil, err
	}

	s.linksCount++
	return link, nil
}

func (s *TestLinksSuite) createUser() (*schema.User, error) {
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

type optionFunc func(*schema.Link)

func withUrl(url string) optionFunc {
	return func(l *schema.Link) {
		l.URL = url
	}
}

func withUserID(id uint64) optionFunc {
	return func(l *schema.Link) {
		l.UserID = id
	}
}
