package api

import (
	"context"
	"database/sql"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/valyala/fasthttp"

	"github.com/ngalayko/url_shortner/server/template"
)

var (
	usersRegEx     = regexp.MustCompile(`/users/(\d+)`)
	userLinksRegEx = regexp.MustCompile(`/users/(\d+)/links`)
)

func (a *Api) getHandlers(appCtx context.Context, requestCtx *fasthttp.RequestCtx) {

	requestUrl := string(requestCtx.RequestURI())
	switch {
	case requestUrl == "/":
		a.renderMainPage(requestCtx)

	case userLinksRegEx.MatchString(requestUrl):
		id, err := parseUserID(requestUrl)
		if err != nil {
			a.responseErr(requestCtx, err)
			return
		}

		a.queryUserLinks(requestCtx, id)

	case usersRegEx.MatchString(requestUrl):
		id, err := parseUserID(requestUrl)
		if err != nil {
			a.responseErr(requestCtx, err)
			return
		}

		a.queryUser(requestCtx, id)

	default:
		a.queryLink(requestCtx)

	}
}

func (a *Api) renderMainPage(ctx *fasthttp.RequestCtx) {
	data, err := template.Index(
		template.WithFacebookConfig(a.fbConfig),
	)
	if err != nil {
		a.responseErr(ctx, err)
		return
	}

	a.responseHtml(ctx, data)
}

func (a *Api) queryLink(ctx *fasthttp.RequestCtx) {
	shortUrl := string(ctx.RequestURI())[1:]

	if len(shortUrl) == 0 {
		ctx.NotFound()
		return
	}

	link, err := a.links.QueryLinkByShortUrl(shortUrl)
	switch {
	case err == sql.ErrNoRows:
		ctx.NotFound()
		return
	case err != nil:
		a.responseErr(ctx, err)
		return
	}

	ctx.Redirect(link.URL, http.StatusFound)
}

func (a *Api) queryUser(ctx *fasthttp.RequestCtx, id uint64) {

	user, err := a.users.QueryUserById(id)
	switch {
	case err == sql.ErrNoRows:
		ctx.NotFound()
		return
	case err != nil:
		a.responseErr(ctx, err)
		return
	}

	a.responseData(ctx, user)
}

func (a *Api) queryUserLinks(ctx *fasthttp.RequestCtx, userID uint64) {

	links, err := a.links.QueryLinksByUser(userID)
	if err != nil {
		a.responseErr(ctx, err)
		return
	}

	a.responseData(ctx, links)
}

func parseUserID(requestURL string) (uint64, error) {
	id := strings.Split(requestURL, "/")[2]

	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, nil
	}

	return uint64(intID), nil
}
