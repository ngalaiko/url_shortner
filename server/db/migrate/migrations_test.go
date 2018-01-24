package migrate

import (
	"context"
	"log"
	"testing"

	. "gopkg.in/check.v1"

	"github.com/ngalayko/url_shortner/server/cache"
	"github.com/ngalayko/url_shortner/server/config"
	"github.com/ngalayko/url_shortner/server/logger"
)

type TestMigrationSuite struct {
	ctx context.Context

	migrations *Migrate
}

func Test(t *testing.T) { TestingT(t) }

var suite *TestMigrationSuite

var _ = Suite(&TestMigrationSuite{})

func (s *TestMigrationSuite) SetUpSuite(c *C) {
	suite = &TestMigrationSuite{
		ctx: context.Background(),
	}
	s.migrations = FromContext(s.ctx)

	m := FromContext(s.ctx)
	if err := m.Flush(); err != nil {
		c.Fatal(err)
	}

	if err := m.Apply(); err != nil {
		log.Panicf("error applying migrations: %s", err)
	}
}

func (s *TestMigrationSuite) init() {
	s.ctx = cache.NewContext(nil, cache.NewStubCache())
	s.ctx = logger.NewContext(s.ctx, logger.NewTestLogger())
	s.ctx = config.NewContext(s.ctx, config.NewTestConfig())
}
