package tables

import (
	"fmt"
	"time"

	. "gopkg.in/check.v1"

	"github.com/ngalayko/url_shortner/server/dao"
	"github.com/ngalayko/url_shortner/server/schema"
)

func (s *TestTablesSuite) Test_InsertLink__should_insert_link(c *C) {
	link := s.testLink(c)

	c.Assert(link.ID, Not(Equals), uint64(0))
}

func (s *TestTablesSuite) Test_GetLinkById__should_select_link(c *C) {
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

func (s *TestTablesSuite) Test_SelectLinksByFields__should_select_links_by_params(c *C) {
	link1 := s.testLink(c)
	link2 := s.testLink(c)
	s.testLink(c)

	param1 := dao.NewParam(1).Add("id", link1.ID)
	param2 := dao.NewParam(1).Add("url", link2.URL)
	params := dao.NewParams(2).Append(param1).Append(param2)

	selected, err := s.service.SelectLinksByFields(params)
	if err != nil {
		c.Fatal(err)
	}

	c.Assert(len(selected), Equals, 2)
	c.Assert(selected[0].ID, Equals, link1.ID)
	c.Assert(selected[1].ID, Equals, link2.ID)
}

func (s *TestTablesSuite) Test_SelectLinksByFields__should_select_link_by_many_params(c *C) {
	link1 := s.testLink(c)

	link2 := s.testLink(c)
	link2.URL = link1.URL
	if err := s.service.UpdateLink(link2); err != nil {
		c.Fatal(err)
	}

	s.testLink(c)

	param1 := dao.NewParam(1).
		Add("id", link1.ID).
		Add("url", link1.URL)
	params := dao.NewParams(2).Append(param1)

	selected, err := s.service.SelectLinksByFields(params)
	if err != nil {
		c.Fatal(err)
	}

	c.Assert(len(selected), Equals, 1)
	c.Assert(selected[0].ID, Equals, link1.ID)
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
