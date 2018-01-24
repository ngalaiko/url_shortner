package server

import (
	"context"

	"go.uber.org/zap"

	"github.com/ngalayko/url_shortner/server/api"
	"github.com/ngalayko/url_shortner/server/cache"
	"github.com/ngalayko/url_shortner/server/config"
	"github.com/ngalayko/url_shortner/server/db"
	"github.com/ngalayko/url_shortner/server/db/migrate"
	"github.com/ngalayko/url_shortner/server/facebook"
	"github.com/ngalayko/url_shortner/server/logger"
	"github.com/ngalayko/url_shortner/server/services/links"
	"github.com/ngalayko/url_shortner/server/services/session"
	"github.com/ngalayko/url_shortner/server/services/token"
	"github.com/ngalayko/url_shortner/server/services/users"
)

// Application is an application main object
type Application struct {
	ctx        context.Context
	cancelFunc context.CancelFunc

	logger logger.ILogger
}

type newServiceFunc func(context.Context, interface{}) context.Context

var (
	services = []newServiceFunc{
		logger.NewContext,
		config.NewContext,
		cache.NewContext,
		db.NewContext,
		migrate.NewContext,
		links.NewContext,
		users.NewContext,
		token.NewContext,
		api.NewContext,
		facebook.NewContext,
		session.NewContext,
	}
)

// NewApplication creates new application
func NewApplication() *Application {

	ctx, cancelFunc := context.WithCancel(context.Background())

	app := &Application{
		ctx:        ctx,
		cancelFunc: cancelFunc,
		logger:     logger.FromContext(ctx),
	}

	for _, service := range services {
		app.ctx = service(app.ctx, nil)
	}

	app.logger = logger.FromContext(app.ctx)
	if err := migrate.FromContext(app.ctx).Apply(); err != nil {
		app.logger.Panic("error while migrations",
			zap.Error(err),
		)
	}

	return app
}

// Serve serve web
func (app *Application) Serve() {
	api.FromContext(app.ctx).Serve()
}
