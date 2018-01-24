package api

import (
	"github.com/valyala/fasthttp"

	"github.com/ngalayko/url_shortner/server/schema"
)

// Ctx wrapper
type Ctx struct {
	*fasthttp.RequestCtx

	Errors []error

	User    *schema.User
	Links   []*schema.Link
	Session *schema.Session

	RedirectUrl string
}

// Authorized true if request from authorized user
func (c *Ctx) Authorized() bool {
	return !(c.User == nil)
}

// AddError adds error to ctx result
func (c *Ctx) AddError(err error) {
	c.Errors = append(c.Errors, err)
}
