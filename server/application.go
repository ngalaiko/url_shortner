package server

import (
	"context"

	"github.com/ngalayko/url_shortner/server/config"
	"github.com/ngalayko/url_shortner/server/db"
	"github.com/ngalayko/url_shortner/server/logger"
)

type Application struct {
	ctx context.Context
}

type newServiceFunc func(context.Context, interface{}) context.Context

var (
	services = []newServiceFunc{
		logger.NewContext,
		config.NewContext,
		db.NewContext,
	}
)

func NewApplication() *Application {

	app := &Application{
		ctx: context.Background(),
	}

	app.initServices()

	return app
}

func (app *Application) initServices() {
	for _, service := range services {
		app.ctx = service(app.ctx, nil)
	}
}

func (app *Application) Healthcheck() {
	l := logger.FromContext(app.ctx)

	l.Info("I'm ok!")
}
