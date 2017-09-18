package server

import (
	"context"

	"go.uber.org/zap"

	"github.com/ngalayko/url_shortner/server/api"
	"github.com/ngalayko/url_shortner/server/dao/migrate"
	"github.com/ngalayko/url_shortner/server/logger"
)

// Application is an application main object
type Application struct {
	ctx        context.Context
	cancelFunc context.CancelFunc

	logger logger.ILogger
}

// NewApplication creates new application
func NewApplication() *Application {

	ctx, cancelFunc := context.WithCancel(context.Background())

	app := &Application{
		ctx:        ctx,
		cancelFunc: cancelFunc,
		logger: logger.FromContext(ctx),
	}

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
