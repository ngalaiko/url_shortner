package web

import (
	"context"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/ngalayko/url_shortner/server/config"
	"github.com/ngalayko/url_shortner/server/logger"
)

const (
	ctxKey webCtxKey = "web_ctx_key"
)

type webCtxKey string

type Web struct {
	handler fasthttp.RequestHandler

	config *config.Config
	logger *logger.Logger
}

func NewContext(ctx context.Context, web interface{}) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if _, ok := web.(*Web); !ok {
		web = newWeb(ctx)
	}

	return context.WithValue(ctx, ctxKey, web)
}

func FromContext(ctx context.Context) *Web {
	if web, ok := ctx.Value(ctxKey).(*Web); ok {
		return web
	}

	return newWeb(ctx)
}

func newWeb(ctx context.Context) *Web {
	w := &Web{
		logger: logger.FromContext(ctx),
		config: config.FromContext(ctx),
	}

	w.initHandler()

	return w
}

func (w *Web) Serve() {
	defer func() {
		recover()
	}()

	w.logger.Info("listening http",
		zap.String("address", w.config.Web.Address),
	)

	if err := fasthttp.ListenAndServe(w.config.Web.Address, w.handler); err != nil {
		w.logger.Error("error while serving",
			zap.Error(err),
		)
	}
}

func (w *Web) initHandler() {
	w.handler = func(ctx *fasthttp.RequestCtx) {
		w.logger.Info("handle request",
			zap.ByteString("method", ctx.Method()),
			zap.ByteString("url", ctx.RequestURI()),
			zap.ByteString("body", ctx.PostBody()),
		)

		switch {
		case ctx.IsGet():
			w.getHandlers(ctx)
		case ctx.IsPost():
			w.postHandlers(ctx)
		default:
			ctx.NotFound()
		}
	}
}

func (w *Web) getHandlers(ctx *fasthttp.RequestCtx) {

	switch string(ctx.RequestURI()) {
	case "/health_check":
		ctx.WriteString("ok")
	default:
		ctx.NotFound()
	}
}

func (w *Web) postHandlers(ctx *fasthttp.RequestCtx) {

	switch string(ctx.RequestURI()) {
	default:
		ctx.NotFound()
	}
}
