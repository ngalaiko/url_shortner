package api

import (
	"context"
	"net/http"

	"github.com/valyala/fasthttp"
)

func (a *Api) getHandlers(appCtx context.Context, requestCtx *fasthttp.RequestCtx) {

	switch string(requestCtx.RequestURI()) {
	case "/health_check":
		a.healthCheck(requestCtx)

	default:
		a.queryLink(requestCtx)

	}
}

func (a *Api) healthCheck(requestCtx *fasthttp.RequestCtx) {

	a.responseData(requestCtx, struct {
		Status bool `json:"status"`
	}{
		Status: true,
	})
}

func (a *Api) queryLink(ctx *fasthttp.RequestCtx) {
	shortUrl := string(ctx.RequestURI())[1:]

	link, err := a.links.QueryLinkByShortUrl(shortUrl)
	if err != nil {
		a.responseErr(ctx, err)
		return
	}

	ctx.Redirect(link.URL, http.StatusFound)
}
