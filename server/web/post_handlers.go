package web

import (
	"context"

	"github.com/valyala/fasthttp"
)

func (w *Web) postHandlers(appCtx context.Context, requestCtx *fasthttp.RequestCtx) {

	switch string(requestCtx.RequestURI()) {
	default:
		requestCtx.NotFound()
	}
}
