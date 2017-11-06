package facebook

import (
	"context"
	"fmt"
	"strings"

	"github.com/huandu/facebook"

	"github.com/ngalayko/url_shortner/server/config"
	"github.com/ngalayko/url_shortner/server/logger"
)

const (
	apiMeURL       = "/me"
	apiFieldsParam = "fields"

	apiFieldsSeparator = ","

	apiID        = "id"
	apiFirstName = "first_name"
	apiLastName  = "last_name"
)

// Api is a facebook api wrapper
type Api struct {
	logger logger.ILogger
	config config.FacebookConfig

	app *facebook.App
}

type User struct {
	ID        string
	FirstName string
	LastName  string
}

func newApi(ctx context.Context) *Api {

	cfg := config.FromContext(ctx).Facebook

	app := facebook.New(cfg.FacebookAppID, cfg.FacebookAppSecret)
	app.RedirectUri = cfg.FacebookLoginURL

	return &Api{
		logger: logger.FromContext(ctx),
		config: cfg,

		app: app,
	}
}

// GetUserByRequest returns user by facebook token from facebook graph api
func (a *Api) GetUserByRequest(facebookCode string) (*User, error) {

	token, err := a.app.ParseCode(facebookCode)
	if err != nil {
		return nil, fmt.Errorf("erorr while parsing code: %s", err)
	}

	session := a.app.Session(token)

	fbResult, err := session.Get(apiMeURL, facebook.Params{
		apiFieldsParam: strings.Join(
			[]string{apiID, apiFirstName, apiLastName},
			apiFieldsSeparator,
		),
	})
	if err != nil {
		return nil, err
	}

	return &User{
		FirstName: fmt.Sprintf("%s", fbResult[apiFirstName]),
		LastName:  fmt.Sprintf("%s", fbResult[apiLastName]),
		ID:        fmt.Sprintf("%s", fbResult[apiID]),
	}, nil
}
