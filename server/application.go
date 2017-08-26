package server

import (
	"context"

	"github.com/ngalayko/url_shortner/server/config"
	"github.com/ngalayko/url_shortner/server/dao"
	"github.com/ngalayko/url_shortner/server/dao/migrate"
	"github.com/ngalayko/url_shortner/server/logger"
	"github.com/ngalayko/url_shortner/server/web"
	"go.uber.org/zap"
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
		dao.NewContext,
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
	web.FromContext(app.ctx).Serve()
}
