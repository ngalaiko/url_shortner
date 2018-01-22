package template

import (
	"bytes"
	"html/template"

	"github.com/ngalayko/url_shortner/server/config"
	"github.com/ngalayko/url_shortner/server/schema"
)

const (
	dataPath         = "template/data/"
	indexFileName    = "index.html"
	notFoundFileName = "not_found.html"
	headFileName     = "head.html"
)

var (
	headTemplate     = template.Must(template.New("head").Parse(string(MustAsset(dataPath + headFileName))))
	indexTemplate    = template.Must(template.New("index").Parse(string(MustAsset(dataPath + indexFileName))))
	notFoundTemplate = template.Must(template.New("notFound").Parse(string(MustAsset(dataPath + notFoundFileName))))
)

type data struct {
	Config config.FacebookConfig
	Errors []error

	User  *schema.User
	Links []*schema.Link
}

// DataFunc is a func to modify template data
type DataFunc func(*data)

// head returns head.html template
func head(d *data) ([]byte, error) {
	var headBuffer bytes.Buffer

	if err := headTemplate.Execute(&headBuffer, d); err != nil {
		return nil, err
	}

	return headBuffer.Bytes(), nil
}

// NotFound return not found page
func NotFound(dataOps ...DataFunc) ([]byte, error) {
	var notFoundBuffer bytes.Buffer

	d := parseOptions(dataOps...)

	headBytes, err := head(d)
	if err != nil {
		return nil, err
	}

	if err := notFoundTemplate.Execute(&notFoundBuffer, d); err != nil {
		return nil, err
	}

	return append(headBytes, notFoundBuffer.Bytes()...), nil
}

// Index return index page
func Index(dataOps ...DataFunc) ([]byte, error) {
	var indexBuffer bytes.Buffer

	d := parseOptions(dataOps...)

	headBytes, err := head(d)
	if err != nil {
		return nil, err
	}

	if err := indexTemplate.Execute(&indexBuffer, d); err != nil {
		return nil, err
	}

	return append(headBytes, indexBuffer.Bytes()...), nil
}

func parseOptions(dataOpts ...DataFunc) *data {
	d := &data{}
	for _, option := range dataOpts {
		option(d)
	}

	return d
}

// WithFacebookConfig sets template facebook config
func WithFacebookConfig(cfg config.FacebookConfig) DataFunc {
	return func(d *data) {
		d.Config = cfg
	}
}

// WithErrors sets template errors
func WithErrors(errors []error) DataFunc {
	return func(d *data) {
		d.Errors = errors
	}
}

// WithUser sets template logged user
func WithUser(user *schema.User) DataFunc {
	return func(d *data) {
		d.User = user
	}
}

// WithLinks sets template links
func WithLinks(links ...*schema.Link) DataFunc {
	return func(d *data) {
		d.Links = links
	}
}
