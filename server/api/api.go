package api

import (
	"context"
	"encoding/json"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/ngalayko/url_shortner/server/config"
	"github.com/ngalayko/url_shortner/server/dao"
	"github.com/ngalayko/url_shortner/server/facebook"
	"github.com/ngalayko/url_shortner/server/logger"
	"github.com/ngalayko/url_shortner/server/services/links"
	"github.com/ngalayko/url_shortner/server/services/user_token"
	"github.com/ngalayko/url_shortner/server/services/users"
)

const (
	ctxKey apiCtxKey = "api_ctx_key"
)

type apiCtxKey string

type response struct {
	Ok   bool        `json:"ok"`
	Data interface{} `json:"data"`
	Err  string      `json:"err"`
}

// Api is a web service
type Api struct {
	handler fasthttp.RequestHandler

	config   config.WebConfig
	fbConfig config.FacebookConfig

	logger logger.ILogger
	db     *dao.Db

	facebookAPI *facebook.Api

	links      *links.Service
	users      *users.Service
	userTokens *user_token.Service
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
	cfg := config.FromContext(ctx)

	w := &Api{
		config:   cfg.Web,
		fbConfig: cfg.Facebook,

		logger: logger.FromContext(ctx),
		db:     dao.FromContext(ctx),

		facebookAPI: facebook.FromContext(ctx),

		links:      links.FromContext(ctx),
		users:      users.FromContext(ctx),
		userTokens: user_token.FromContext(ctx),
	}

	w.initHandler(ctx)

	return w
}

// Serve serve web with config credentials
func (a *Api) Serve() {
	defer func() {
		recover()
	}()

	go func() {
		a.logger.Info("listening pprof",
			zap.String("address", a.config.PprofAddress),
		)
		if err := http.ListenAndServe(a.config.PprofAddress, nil); err != nil {
			a.logger.Error("error while start pprof",
				zap.Error(err),
			)
		}
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

		ctx, err := a.NewCtx(requestCtx)
		if err != nil {
			a.responseErr(ctx, err)
			return
		}

		switch {
		case requestCtx.IsGet():
			a.getHandlers(appCtx, ctx)

		case requestCtx.IsPost():
			a.postHandlers(appCtx, ctx)

		default:
			requestCtx.NotFound()

		}

		if ctx.RedirectUrl != "" {
			ctx.Redirect(ctx.RedirectUrl, http.StatusFound)
		}

		a.logger.Info("handle request",
			zap.ByteString("method", ctx.Method()),
			zap.ByteString("url", ctx.RequestURI()),
			zap.ByteString("body", ctx.PostBody()),
			zap.Bool("authorized", ctx.Authorized()),
			zap.Duration("duration", time.Since(start)),
		)
	}
}

func (a *Api) responseErr(ctx *Ctx, err error) {
	ctx.Response.SetStatusCode(http.StatusBadRequest)
	ctx.Response.Header.Set("Content-Type", "application/json")

	data, err := json.Marshal(response{
		Ok:  false,
		Err: err.Error(),
	})
	if err != nil {
		a.responseErr(ctx, err)
	}

	a.responseBytes(ctx, data)
}

func (a *Api) responseData(ctx *Ctx, obj interface{}) {
	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.Header.Set("Content-Type", "application/json")

	data, err := json.Marshal(response{
		Ok:   true,
		Data: obj,
	})
	if err != nil {
		a.responseErr(ctx, err)
	}

	a.responseBytes(ctx, data)
}

func (a *Api) responseHtml(ctx *Ctx, data []byte) {
	ctx.Response.Header.Set("Content-Type", "text/html; charset=utf-8")

	a.responseBytes(ctx, data)
}

func (a *Api) responseBytes(ctx *Ctx, data []byte) {
	ctx.Response.AppendBody(data)
}
