package api

import (
	"context"
	"time"

	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"

	"github.com/ngalayko/url_shortner/server/schema"
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
	link := &schema.Link{}
	if err := easyjson.Unmarshal(ctx.PostBody(), link); err != nil {
		a.responseErr(ctx, err)
		return
	}

	link.CreatedAt = time.Now()
	if err := a.tables.InsertLink(link); err != nil {
		a.responseErr(ctx, err)
		return
	}

	a.responseData(ctx, link)
}
