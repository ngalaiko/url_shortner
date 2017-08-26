package web

import "github.com/valyala/fasthttp"

func (w *Web) getHandlers(ctx *fasthttp.RequestCtx) {

	switch string(ctx.RequestURI()) {
	case "/health_check":
		ctx.WriteString("ok")
	default:
		ctx.NotFound()
	}
}
