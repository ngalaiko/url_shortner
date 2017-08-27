package tables

import (
	"context"
	"log"
	"testing"

	"go.uber.org/zap"
	. "gopkg.in/check.v1"

	"github.com/ngalayko/url_shortner/server/config"
	"github.com/ngalayko/url_shortner/server/dao/migrate"
	"github.com/ngalayko/url_shortner/server/logger"
)

type TestTablesSuite struct {
	ctx context.Context

	service *Tables

	usersCount int
	linksCount int
}

var suite *TestTablesSuite

func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&TestTablesSuite{})

func (s *TestTablesSuite) SetUpSuite(c *C) {
	suite = &TestTablesSuite{
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

func (s *TestTablesSuite) init() {
	s.ctx = logger.NewContext(s.ctx, s.initLogger())
	s.ctx = config.NewContext(s.ctx, s.initConfig())
	s.ctx = migrate.NewContext(s.ctx, nil)

	s.service = FromContext(s.ctx)
}

func (s *TestTablesSuite) initLogger() *logger.Logger {
	l, err := zap.NewDevelopment()
	if err != nil {
		log.Panicf("error while init logger: %s ", err)
	}

	return &logger.Logger{l}
}

func (s *TestTablesSuite) initConfig() *config.Config {
	return &config.Config{
		Db: config.DbConfig{
			Driver:       "postgres",
			Connect:      "host=localhost user=url_short_test dbname=url_short_test sslmode=disable password=secret",
			MaxIdleConns: 5,
			MaxOpenConns: 5,
		},
	}
}
