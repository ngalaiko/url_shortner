package api

import (
	"context"

	"github.com/valyala/fasthttp"

	"github.com/ngalayko/url_shortner/client/template"
)

func (a *Api) getHandlers(appCtx context.Context, requestCtx *fasthttp.RequestCtx) {

	requestUrl := string(requestCtx.RequestURI())
	switch requestUrl {

	case "/":
		a.renderIndexPage(requestCtx)

	default:
		requestCtx.NotFound()

	}
}

func (a *Api) renderIndexPage(ctx *fasthttp.RequestCtx) {

	data, err := template.Index()
	if err != nil {
		a.responseErr(ctx, err)
	}

	a.responseBytes(ctx, data)
}
