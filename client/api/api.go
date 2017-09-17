package api

import (
	"context"
	"encoding/json"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/ngalayko/url_shortner/client/config"
	"github.com/ngalayko/url_shortner/client/logger"
)

const (
	ctxKey apiCtxKey = "api_ctx_key"
)

type apiCtxKey string

type errResponse struct {
	Err string `json:"err"`
}

// Api is a web service
type Api struct {
	handler fasthttp.RequestHandler
	config  config.WebConfig
	logger  *logger.Logger
}

// NewContext stores web in context
func NewContext(ctx context.Context, web interface{}) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if _, ok := web.(*Api); !ok {
		web = newApi(ctx)
	}

	return context.WithValue(ctx, ctxKey, web)
}

// FromContext return web from context
func FromContext(ctx context.Context) *Api {
	if web, ok := ctx.Value(ctxKey).(*Api); ok {
		return web
	}

	return newApi(ctx)
}

func newApi(ctx context.Context) *Api {
	w := &Api{
		config: config.FromContext(ctx).Web,
		logger: logger.FromContext(ctx).Prefix("api"),
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
		start := time.Now()

		switch {
		case requestCtx.IsGet():
			a.getHandlers(appCtx, requestCtx)
		default:
			requestCtx.NotFound()
		}

		a.logger.Info("handle request",
			zap.ByteString("method", requestCtx.Method()),
			zap.ByteString("url", requestCtx.RequestURI()),
			zap.ByteString("body", requestCtx.PostBody()),
			zap.Duration("duration", time.Since(start)),
		)
	}
}

func (a *Api) responseErr(ctx *fasthttp.RequestCtx, err error) {
	data, err := json.Marshal(errResponse{
		Err: err.Error(),
	})
	if err != nil {
		a.responseErr(ctx, err)
	}

	ctx.Response.SetStatusCode(http.StatusInternalServerError)
	ctx.Response.AppendBody(data)
}

func (a *Api) responseBytes(ctx *fasthttp.RequestCtx, data []byte) {
	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.AppendBody(data)
}
