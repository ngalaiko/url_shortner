package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	// activate pprof
	_ "net/http/pprof"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/ngalayko/url_shortner/server/config"
	"github.com/ngalayko/url_shortner/server/db"
	"github.com/ngalayko/url_shortner/server/facebook"
	"github.com/ngalayko/url_shortner/server/logger"
	"github.com/ngalayko/url_shortner/server/services/links"
	"github.com/ngalayko/url_shortner/server/services/session"
	"github.com/ngalayko/url_shortner/server/services/token"
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

// API is a web service
type API struct {
	handler fasthttp.RequestHandler

	config   config.WebConfig
	fbConfig config.FacebookConfig

	logger logger.ILogger
	db     *db.Db

	facebookAPI *facebook.API

	links      *links.Service
	users      *users.Service
	userTokens *token.Service
	sessions   session.ISession
}

// NewContext stores web in context
func NewContext(ctx context.Context, web interface{}) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if _, ok := web.(*API); !ok {
		web = newAPI(ctx)
	}

	return context.WithValue(ctx, ctxKey, web)
}

// FromContext return web from context
func FromContext(ctx context.Context) *API {
	if web, ok := ctx.Value(ctxKey).(*API); ok {
		return web
	}

	return newAPI(ctx)
}

func newAPI(ctx context.Context) *API {
	cfg := config.FromContext(ctx)

	w := &API{
		config:   cfg.Web,
		fbConfig: cfg.Facebook,

		logger: logger.FromContext(ctx),
		db:     db.FromContext(ctx),

		facebookAPI: facebook.FromContext(ctx),

		links:      links.FromContext(ctx),
		users:      users.FromContext(ctx),
		userTokens: token.FromContext(ctx),
		sessions:   session.FromContext(ctx),
	}

	w.initHandler(ctx)

	return w
}

// Serve serve web with config credentials
func (a *API) Serve() {
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

func (a *API) initHandler(appCtx context.Context) {
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
			a.responseNotFound(ctx)

		}

		if ctx.RedirectURL != "" {
			ctx.Redirect(ctx.RedirectURL, http.StatusFound)
		}

		a.logger.Info("handle request",
			zap.ByteString("method", ctx.Method()),
			zap.ByteString("url", ctx.RequestURI()),
			zap.ByteString("body", ctx.PostBody()),
			zap.Reflect("user", ctx.User),
			zap.Duration("duration", time.Since(start)),
			zap.Reflect("session", ctx.Session),
		)
	}
}

func (a *API) responseErr(ctx *Ctx, err error) {
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

func (a *API) responseData(ctx *Ctx, obj interface{}) {
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

func (a *API) responseHTML(ctx *Ctx, data []byte) {
	ctx.Response.Header.Set("Content-Type", "text/html; charset=utf-8")

	a.responseBytes(ctx, data)
}

func (a *API) responseBytes(ctx *Ctx, data []byte) {
	ctx.Response.AppendBody(data)
}

func (a *API) responseNotFound(ctx *Ctx) {
	data, err := a.renderNotFoundPage(ctx)
	if err != nil {
		a.responseErr(ctx, err)
		return
	}

	a.responseHTML(ctx, data)
}
