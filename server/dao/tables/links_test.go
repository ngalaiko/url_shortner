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

	selected, err := s.service.SelectLinkById(link.ID)
	if err != nil {
		c.Fatal(err)
	}

	c.Assert(link, DeepEquals, selected)
}

func (s *TestTablesSuite) Test_UpdateLink__should_update_link(c *C) {
	link := s.testLink(c)
	link.URL = "updated"

	if err := s.service.UpdateLink(link); err != nil {
		c.Fatal(err)
	}

	updated, err := s.service.SelectLinkById(link.ID)
	if err != nil {
		c.Fatal(err)
	}

	c.Assert(link, DeepEquals, updated)
}

func (s *TestTablesSuite) testLink(c *C) *schema.Link {
	s.linksCount++

	link := &schema.Link{
		UserID:    s.testUser(c).ID,
		URL:       fmt.Sprintf("link %d", s.linksCount),
		CreatedAt: time.Now(),
	}

	if err := s.service.InsertLink(link); err != nil {
		c.Fatal(err)
	}

	return link
}
