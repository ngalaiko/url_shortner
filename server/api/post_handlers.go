package api

import (
	"context"
	"encoding/json"

	"github.com/ngalayko/url_shortner/server/schema"
)

func (a *Api) postHandlers(appCtx context.Context, requestCtx *Ctx) {

	switch string(requestCtx.RequestURI()) {
	case "/link":
		a.createLink(requestCtx)

	default:
		requestCtx.NotFound()

	}
}

func (a *Api) createLink(ctx *Ctx) {
	link := &schema.Link{}
	if err := json.Unmarshal(ctx.PostBody(), link); err != nil {
		a.responseErr(ctx, err)
		return
	}

	if ctx.Authorized() {
		link.UserID = ctx.User.ID
	}

	if err := a.links.CreateLink(link); err != nil {
		a.responseErr(ctx, err)
		return
	}

	if !ctx.Authorized() {
		ctx.Session.LinkIDs = append(ctx.Session.LinkIDs, link.ID)
		if err := a.sessions.Update(ctx.Session); err != nil {
			a.responseErr(ctx, err)
			return
		}
	}

	a.responseData(ctx, link)
}
