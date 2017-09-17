package client

import (
	"context"

	"github.com/ngalayko/url_shortner/client/api"
)

// Application is a main struct
type Application struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
}

// NewApplication is an application constructor
func NewApplication() *Application {
	ctx, cancelFunc := context.WithCancel(context.Background())

	return &Application{
		ctx:        ctx,
		cancelFunc: cancelFunc,
	}
}

// Serve serve clients
func (a *Application) Serve() {
	api.FromContext(a.ctx).Serve()
}
