package api

import (
	"context"

	"github.com/valyala/fasthttp"
)

func (a *Api) postHandlers(appCtx context.Context, requestCtx *fasthttp.RequestCtx) {

	switch string(requestCtx.RequestURI()) {
	default:
		requestCtx.NotFound()
	}
}
