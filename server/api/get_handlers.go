package api

import (
	"context"
	"database/sql"
	"strings"

	"github.com/ngalayko/url_shortner/server/template"
)

const (
	facebookAccessCodeURLParam = "code"
)

func (a *Api) getHandlers(appCtx context.Context, ctx *Ctx) {

	requestUrl := string(ctx.RequestURI())
	switch {
	case requestUrl == "/":
		data, err := a.renderMainPage(ctx)
		if err != nil {
			ctx.AddError(err)
		}

		a.responseHtml(ctx, data)

	case requestUrl == "/logout":
		if err := a.deleteUserCookie(ctx); err != nil {
			a.responseErr(ctx, err)
			return
		}

	case strings.HasPrefix(requestUrl, facebookLoginRequestURI):
		user, err := a.authorizeUser(ctx.RequestCtx)
		if err != nil {
			a.responseErr(ctx, err)
			return
		}

		if err := a.setUserCookie(ctx.RequestCtx, user); err != nil {
			a.responseErr(ctx, err)
			return
		}

		a.redirectHome(ctx)

	default:
		a.redirectLink(ctx)

	}
}

func (a *Api) redirectHome(ctx *Ctx) {
	ctx.RedirectUrl = "https://" + string(ctx.URI().Host())
}

func (a *Api) renderNotFoundPage(ctx *Ctx) ([]byte, error) {
	data, err := template.NotFound(
		template.WithFacebookConfig(a.fbConfig),
		template.WithUser(ctx.User),
		template.WithErrors(ctx.Errors),
		template.WithLinks(ctx.Links...),
	)
	if err != nil {
		return nil, err
	}

	return data, err
}

func (a *Api) renderMainPage(ctx *Ctx) ([]byte, error) {
	data, err := template.Index(
		template.WithFacebookConfig(a.fbConfig),
		template.WithUser(ctx.User),
		template.WithErrors(ctx.Errors),
		template.WithLinks(ctx.Links...),
	)
	if err != nil {
		return nil, err
	}

	return data, err
}

func (a *Api) redirectLink(ctx *Ctx) {
	shortUrl := string(ctx.RequestURI())[1:]

	if len(shortUrl) == 0 {
		a.responseNotFound(ctx)
		return
	}

	link, err := a.links.QueryLinkByShortUrl(shortUrl)
	switch {
	case err == sql.ErrNoRows:
		a.responseNotFound(ctx)
		return

	case err != nil:
		a.responseErr(ctx, err)
		return

	}

	ctx.RedirectUrl = link.URL
}
