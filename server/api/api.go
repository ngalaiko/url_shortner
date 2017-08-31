package api

import (
	"context"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/ngalayko/url_shortner/server/config"
	"github.com/ngalayko/url_shortner/server/logger"
	"github.com/ngalayko/url_shortner/server/dao/tables"
)

const (
	ctxKey apiCtxKey = "api_ctx_key"
)

type apiCtxKey string

// Api is a web service
type Api struct {
	handler fasthttp.RequestHandler
	config config.WebConfig

	tables *tables.Tables
	logger *logger.Logger
}

// NewContext stores web in context
func NewContext(ctx context.Context, web interface{}) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if _, ok := web.(*Api); !ok {
		web = newWeb(ctx)
	}

	return context.WithValue(ctx, ctxKey, web)
}

// FromContext return web from context
func FromContext(ctx context.Context) *Api {
	if web, ok := ctx.Value(ctxKey).(*Api); ok {
		return web
	}

	return newWeb(ctx)
}

func newWeb(ctx context.Context) *Api {
	w := &Api{
		config: config.FromContext(ctx).Web,

		logger: logger.FromContext(ctx),
		tables: tables.FromContext(ctx),
	}

	w.initHandler(ctx)

	return w
}

// Serve serve web with config credentials
func (a *Api) Serve() {
	defer func() {
		recover()
	}()

	a.logger.Info("listening http",
		zap.String("address", a.config.Address),
	)

	if err := fasthttp.ListenAndServe(a.config.Address, a.handler); err != nil {
		a.logger.Error("error while serving",
			zap.Error(err),
		)
	}
}

func (a *Api) initHandler(appCtx context.Context) {
	a.handler = func(requestCtx *fasthttp.RequestCtx) {
		a.logger.Info("handle request",
			zap.ByteString("method", requestCtx.Method()),
			zap.ByteString("url", requestCtx.RequestURI()),
			zap.ByteString("body", requestCtx.PostBody()),
		)

		switch {
		case requestCtx.IsGet():
			a.getHandlers(appCtx, requestCtx)
		case requestCtx.IsPost():
			a.postHandlers(appCtx, requestCtx)
		default:
			requestCtx.NotFound()
		}
	}
}
