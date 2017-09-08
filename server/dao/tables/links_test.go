package tables

import (
	"fmt"
	"time"

	. "gopkg.in/check.v1"

	"github.com/ngalayko/url_shortner/server/schema"
)

func (s *TestTablesSuite) Test_InsertLink__should_insert_link(c *C) {
	link := s.testLink(c)

	c.Assert(link.ID, Not(Equals), uint64(0))
}

func (s *TestTablesSuite) Test_SelectLink__should_select_link(c *C) {
	link := s.testLink(c)

	selected, err := s.service.GetLinkById(link.ID)
	if err != nil {
		c.Fatal(err)
	}

	c.Assert(link.ID, Equals, selected.ID)
	c.Assert(link.URL, Equals, selected.URL)
	c.Assert(link.ShortURL, Equals, selected.ShortURL)
	c.Assert(link.Views, Equals, selected.Views)
}

func (s *TestTablesSuite) Test_UpdateLink__should_update_link(c *C) {
	link := s.testLink(c)
	link.URL = "updated"

	if err := s.service.UpdateLink(link); err != nil {
		c.Fatal(err)
	}

	updated, err := s.service.GetLinkById(link.ID)
	if err != nil {
		c.Fatal(err)
	}

	c.Assert(link.ID, Equals, updated.ID)
	c.Assert(link.URL, Equals, updated.URL)
	c.Assert(link.ShortURL, Equals, updated.ShortURL)
	c.Assert(link.Views, Equals, updated.Views)
}

func (s *TestTablesSuite) testLink(c *C) *schema.Link {
	s.linksCount++

	link := &schema.Link{
		UserID:    s.testUser(c).ID,
		URL:       fmt.Sprintf("link url %d", s.linksCount),
		ShortURL:  fmt.Sprintf("link short url %d", s.linksCount),
		CreatedAt: time.Now(),
	}

	if err := s.service.InsertLink(link); err != nil {
		c.Fatal(err)
	}

	return link
}
