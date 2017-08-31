package server

import (
	"context"

	"go.uber.org/zap"

	"github.com/ngalayko/url_shortner/server/api"
	"github.com/ngalayko/url_shortner/server/cache"
	"github.com/ngalayko/url_shortner/server/config"
	"github.com/ngalayko/url_shortner/server/dao"
	"github.com/ngalayko/url_shortner/server/dao/migrate"
	"github.com/ngalayko/url_shortner/server/dao/tables"
	"github.com/ngalayko/url_shortner/server/logger"
)

// Application is an application main object
type Application struct {
	ctx context.Context
}

type newServiceFunc func(context.Context, interface{}) context.Context

var (
	services = []newServiceFunc{
		logger.NewContext,
		config.NewContext,
		cache.NewContext,
		dao.NewContext,
		tables.NewContext,
		migrate.NewContext,
	}
)

// NewApplication creates new application
func NewApplication() *Application {

	app := &Application{
		ctx: context.Background(),
	}

	app.initServices()

	l := logger.FromContext(app.ctx)
	if err := migrate.FromContext(app.ctx).Apply(); err != nil {
		l.Panic("error while migrations",
			zap.Error(err),
		)
	}

	return app
}

func (app *Application) initServices() {
	for _, service := range services {
		app.ctx = service(app.ctx, nil)
	}
}

// Serve serve web
func (app *Application) Serve() {
	api.FromContext(app.ctx).Serve()
}
