package server

import (
	"context"

	"go.uber.org/zap"

	"github.com/ngalayko/url_shortner/server/api"
	"github.com/ngalayko/url_shortner/server/cache"
	"github.com/ngalayko/url_shortner/server/config"
	"github.com/ngalayko/url_shortner/server/dao"
	"github.com/ngalayko/url_shortner/server/dao/migrate"
	"github.com/ngalayko/url_shortner/server/facebook"
	"github.com/ngalayko/url_shortner/server/logger"
	"github.com/ngalayko/url_shortner/server/services/links"
	"github.com/ngalayko/url_shortner/server/services/user_token"
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
		dao.NewContext,
		migrate.NewContext,
		links.NewContext,
		users.NewContext,
		user_token.NewContext,
		api.NewContext,
		facebook.NewContext,
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
