package client

import "context"

// Application is a main struct
type Application struct {
	ctx context.Context
}

// NewApplication is an application constructor
func NewApplication() *Application {
	return &Application{
		ctx: context.Background(),
	}
}

// Serve serve clients
func (a *Application) Serve() {

}
