package api

import (
	"context"

	"github.com/valyala/fasthttp"
)

func (a *Api) getHandlers(appCtx context.Context, requestCtx *fasthttp.RequestCtx) {

	switch string(requestCtx.RequestURI()) {
	case "/health_check":
		requestCtx.WriteString("ok")
	default:
		requestCtx.NotFound()
	}
}
