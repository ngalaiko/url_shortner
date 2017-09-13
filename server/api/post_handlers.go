package api

import (
	"context"
	"encoding/json"

	"github.com/valyala/fasthttp"

	"github.com/ngalayko/url_shortner/server/schema"
)

func (a *Api) postHandlers(appCtx context.Context, requestCtx *fasthttp.RequestCtx) {

	switch string(requestCtx.RequestURI()) {
	case "/link":
		a.createLink(requestCtx)
	default:
		requestCtx.NotFound()
	}
}

func (a *Api) createLink(ctx *fasthttp.RequestCtx) {
	link := &schema.Link{}
	if err := json.Unmarshal(ctx.PostBody(), link); err != nil {
		a.responseErr(ctx, err)
		return
	}

	if err := a.links.CreateLink(link); err != nil {
		a.responseErr(ctx, err)
		return
	}

	a.responseData(ctx, link)
}
