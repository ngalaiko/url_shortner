package api

import (
	"database/sql"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/ngalayko/url_shortner/server/schema"
	"github.com/ngalayko/url_shortner/server/services/session"
)

const (
	facebookLoginRequestURI = "/login/facebook"
	userTokenCookie         = "token"
	sessionCookieName       = "session"
)

// NewCtx from request ctx
func (a *Api) NewCtx(requestCtx *fasthttp.RequestCtx) (*Ctx, error) {
	ctx := &Ctx{
		RequestCtx: requestCtx,
	}

	var err error
	ctx.Session, err = a.getSession(requestCtx)
	if err != nil {
		return ctx, err
	}

	user, err := a.getUserFromCookie(requestCtx)
	switch {
	case err == sql.ErrNoRows:
	case err != nil:
		return ctx, err

	default:
		ctx.Response.Header.DelCookie(sessionCookieName)

		ctx.User = user

		ctx.Links, err = a.links.QueryLinksByUser(user.ID)
		if err != nil {
			ctx.AddError(err)
		}

		if len(ctx.Session.LinkIDs) == 0 {
			return ctx, nil
		}

		if err := a.links.TransferLinks(user.ID, ctx.Session.LinkIDs...); err != nil {
			return ctx, err
		}

		ctx.Session.LinkIDs = []uint64{}
		if err := a.sessions.Update(ctx.Session); err != nil {
			return ctx, err
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

func (a *Api) getSession(ctx *fasthttp.RequestCtx) (*schema.Session, error) {
	sessionKey := string(ctx.Request.Header.Cookie(sessionCookieName))

	s, err := a.sessions.Load(sessionKey)
	switch {
	case err == session.ErrorNoSuchSession:
		s = a.sessions.Create()

	case err != nil:
		return nil, err
	}

	sessionCookie := fasthttp.AcquireCookie()
	sessionCookie.SetKey(sessionCookieName)
	sessionCookie.SetValue(s.Key)
	sessionCookie.SetDomainBytes(ctx.URI().Host())
	sessionCookie.SetSecure(true)
	sessionCookie.SetPath("/")
	sessionCookie.SetHTTPOnly(true)
	ctx.Response.Header.SetCookie(sessionCookie)
	fasthttp.ReleaseCookie(sessionCookie)

	return s, nil
}

func (a *Api) getUserFromCookie(ctx *fasthttp.RequestCtx) (*schema.User, error) {
	token := ctx.Request.Header.Cookie(userTokenCookie)

	userToken, err := a.userTokens.GetUserToken(string(token))
	if err != nil {
		return nil, err
	}

	a.logger.Info("user token from cookie",
		zap.Reflect("user token", userToken),
		zap.ByteString("cookie", token),
	)
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
