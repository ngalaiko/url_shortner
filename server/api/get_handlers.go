package api

import (
	"context"
	"database/sql"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/valyala/fasthttp"
)

var (
	usersRegEx = regexp.MustCompile(`/users/(\d+)`)
)

func (a *Api) getHandlers(appCtx context.Context, requestCtx *fasthttp.RequestCtx) {

	requestUrl := string(requestCtx.RequestURI())
	switch {
	case usersRegEx.MatchString(requestUrl):
		a.queryUser(requestCtx, strings.Split(requestUrl, "/")[2])

	default:
		a.queryLink(requestCtx)

	}
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

func (a *Api) queryUser(ctx *fasthttp.RequestCtx, id string) {

	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		a.responseErr(ctx, err)
		return
	}

	user, err := a.users.QueryUserById(uint64(intID))
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
