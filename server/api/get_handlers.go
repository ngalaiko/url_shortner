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

func (a *API) getHandlers(appCtx context.Context, ctx *Ctx) {

	requestURL := string(ctx.RequestURI())
	switch {
	case requestURL == "/":
		data, err := a.renderMainPage(ctx)
		if err != nil {
			ctx.AddError(err)
		}

		a.responseHTML(ctx, data)

	case requestURL == "/logout":
		if err := a.deleteUserCookie(ctx); err != nil {
			a.responseErr(ctx, err)
			return
		}

	case strings.HasPrefix(requestURL, facebookLoginRequestURI):
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

func (a *API) redirectHome(ctx *Ctx) {
	ctx.RedirectURL = "https://" + string(ctx.URI().Host())
}

func (a *API) renderNotFoundPage(ctx *Ctx) ([]byte, error) {
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

func (a *API) renderMainPage(ctx *Ctx) ([]byte, error) {
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

func (a *API) redirectLink(ctx *Ctx) {
	shortURL := string(ctx.RequestURI())[1:]

	if len(shortURL) == 0 {
		a.responseNotFound(ctx)
		return
	}

	link, err := a.links.QueryLinkByShortURL(shortURL)
	switch {
	case err == sql.ErrNoRows:
		a.responseNotFound(ctx)
		return

	case err != nil:
		a.responseErr(ctx, err)
		return

	}

	ctx.RedirectURL = link.URL
}
