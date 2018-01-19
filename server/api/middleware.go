package api

import (
	"database/sql"

	"github.com/valyala/fasthttp"

	"github.com/ngalayko/url_shortner/server/schema"
)

const (
	facebookLoginRequestURI = "/login/facebook"
)

// NewCtx from request ctx
func (a *Api) NewCtx(requestCtx *fasthttp.RequestCtx) (*Ctx, error) {
	ctx := &Ctx{
		RequestCtx: requestCtx,
	}

	user, err := a.getUserFromCookie(requestCtx)
	switch {
	case err == sql.ErrNoRows:

	case err != nil:
		return ctx, err

	default:
		ctx.User = user

		ctx.Links, err = a.links.QueryLinksByUser(user.ID)
		if err != nil {
			ctx.AddError(err)
		}
	}

	return ctx, nil
}

func (a *Api) authorizeUser(ctx *fasthttp.RequestCtx) (*schema.User, error) {
	code := parseFacebookCode(ctx)

	if len(code) == 0 {
		return nil, nil
	}

	facebookUser, err := a.facebookAPI.GetUserByRequest(code)
	if err != nil {
		return nil, err
	}

	return a.users.QueryUserByFacebookUser(facebookUser)
}

func (a *Api) getUserFromCookie(ctx *fasthttp.RequestCtx) (*schema.User, error) {

	token := ctx.Request.Header.Cookie(userTokenCookie)

	userToken, err := a.userTokens.GetUserToken(string(token))
	if err != nil {
		return nil, err
	}

	return a.users.QueryUserById(userToken.UserID)
}

func (a *Api) deleteUserCookie(ctx *Ctx) error {
	if !ctx.Authorized() {
		return nil
	}

	token := ctx.Request.Header.Cookie(userTokenCookie)
	if err := a.userTokens.DeleteUserToken(ctx.User.ID, string(token)); err != nil {
		return err
	}
	ctx.Response.Header.DelClientCookie(userTokenCookie)

	return nil
}

func (a *Api) setUserCookie(ctx *fasthttp.RequestCtx, user *schema.User) error {

	userToken, err := a.userTokens.CreateUserToken(user)
	if err != nil {
		return err
	}

	tokenCookie := fasthttp.AcquireCookie()
	tokenCookie.SetKey(userTokenCookie)
	tokenCookie.SetValue(userToken.Token)
	tokenCookie.SetExpire(userToken.ExpiredAt)
	tokenCookie.SetDomainBytes(ctx.URI().Host())
	tokenCookie.SetSecure(true)
	tokenCookie.SetPath("/")
	tokenCookie.SetHTTPOnly(true)
	ctx.Response.Header.SetCookie(tokenCookie)
	fasthttp.ReleaseCookie(tokenCookie)

	return nil
}

func parseFacebookCode(ctx *fasthttp.RequestCtx) string {
	return string(ctx.QueryArgs().Peek(facebookAccessCodeURLParam))
}
