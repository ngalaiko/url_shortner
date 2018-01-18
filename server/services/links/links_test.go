package links

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	. "gopkg.in/check.v1"

	"github.com/ngalayko/url_shortner/server/cache"
	"github.com/ngalayko/url_shortner/server/config"
	"github.com/ngalayko/url_shortner/server/dao/migrate"
	"github.com/ngalayko/url_shortner/server/logger"
	"github.com/ngalayko/url_shortner/server/schema"
)

type TestLinksSuite struct {
	ctx context.Context

	service *Service

	usersCount int
	linksCount int
}

func Test(t *testing.T) { TestingT(t) }

var suite *TestLinksSuite

var _ = Suite(&TestLinksSuite{})

func (s *TestLinksSuite) SetUpSuite(c *C) {
	suite = &TestLinksSuite{
		ctx: context.Background(),
	}

	s.init()

	m := migrate.FromContext(s.ctx)
	if err := m.Flush(); err != nil {
		c.Fatal(err)
	}

	if err := m.Apply(); err != nil {
		log.Panicf("error applying migrations: %s", err)
	}
}

func (s *TestLinksSuite) init() {
	s.ctx = cache.NewContext(nil, cache.NewStubCache())
	s.ctx = logger.NewContext(s.ctx, logger.NewTestLogger())
	s.ctx = config.NewContext(s.ctx, config.NewTestConfig())
	s.ctx = migrate.NewContext(s.ctx, nil)

	s.service = FromContext(s.ctx)
}

func (s *TestLinksSuite) Test_CreateLink__should_create_link(c *C) {
	_, err := s.createLink()
	c.Assert(err, IsNil)
}

func (s *TestLinksSuite) Test_CreateLink__should_return_existing_link(c *C) {
	link1, err := s.createLink()
	if err != nil {
		c.Fatal(err)
	}

	link2, err := s.createLink(
		withUrl(link1.URL),
	)
	if err != nil {
		c.Fatal(err)
	}

	c.Assert(link1.ID, Equals, link2.ID)
}

func (s *TestLinksSuite) Test_CreateLink__should_return_link_if_not_valid_exists(c *C) {
	link1, err := s.createLink()
	if err != nil {
		c.Fatal(err)
	}

	if err := s.service.deleteLink(link1); err != nil {
		c.Fatal(err)
	}

	link2, err := s.createLink(
		withUrl(link1.URL),
	)
	if err != nil {
		c.Fatal(err)
	}

	link3, err := s.createLink(
		withUrl(link1.URL),
	)
	if err != nil {
		c.Fatal(err)
	}

	c.Assert(link1.ID, Not(Equals), link2.ID)
	c.Assert(link3.ID, Equals, link2.ID)
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

func (s *TestLinksSuite) Test_prepareLink__should_set_created_at_expired_at(c *C) {

	link := &schema.Link{
		URL: "vk.com",
	}

	if err := prepareLink(link); err != nil {
		c.Fatal(err)
	}

	c.Assert(link.CreatedAt, Not(Equals), time.Unix(0, 0))
	c.Assert(link.ExpiredAt.After(link.CreatedAt), Equals, true)
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

// helpers

func (s *TestLinksSuite) createLink(opts ...optionFunc) (*schema.Link, error) {

	link := &schema.Link{
		URL: fmt.Sprintf("http://vk.com/%d", s.linksCount),
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

type optionFunc func(*schema.Link)

func withUrl(url string) optionFunc {
	return func(l *schema.Link) {
		l.URL = url
	}
}
