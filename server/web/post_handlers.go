package web

func (w *Web) postHandlers(ctx *fasthttp.RequestCtx) {

	switch string(ctx.RequestURI()) {
	default:
		ctx.NotFound()
	}
}
