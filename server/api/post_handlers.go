package api

import (
	"context"

	"github.com/ngalayko/url_shortner/server/schema"
	"github.com/valyala/fasthttp"
)

func (a *Api) postHandlers(appCtx context.Context, requestCtx *fasthttp.RequestCtx) {

	switch string(requestCtx.RequestURI()) {
	case "/create_link":
		a.createLink(requestCtx)
	default:
		requestCtx.NotFound()
	}
}

func (a *Api) createLink(ctx *fasthttp.RequestCtx) {

}
