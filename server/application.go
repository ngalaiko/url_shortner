package server

import (
	"context"

	"github.com/ngalayko/url_shortner/server/config"
	"github.com/ngalayko/url_shortner/server/dao"
	"github.com/ngalayko/url_shortner/server/dao/migrate"
	"github.com/ngalayko/url_shortner/server/logger"
	"go.uber.org/zap"
)

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

func (app *Application) Healthcheck() {
	l := logger.FromContext(app.ctx)

	l.Info("I'm ok!")
}
