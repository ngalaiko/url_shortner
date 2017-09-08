package links

import (
	"testing"
	"time"

	. "gopkg.in/check.v1"

	"github.com/ngalayko/url_shortner/server/schema"
)

type TestLinksSuite struct{}

func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&TestLinksSuite{})

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
