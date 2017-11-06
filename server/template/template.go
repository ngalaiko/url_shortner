package template

import (
	"bytes"
	"html/template"

	"github.com/ngalayko/url_shortner/server/config"
	"github.com/ngalayko/url_shortner/server/schema"
)

const (
	dataPath = "template/data/"
	index    = "index.html"
)

var (
	indexTemplate = template.Must(template.New("index").Parse(string(MustAsset(dataPath + index))))
)

type data struct {
	Config config.FacebookConfig
	Errors []error

	User  *schema.User
	Links []*schema.Link
}

// DataFunc is a func to modify template data
type DataFunc func(*data)

// Index return index.html template
func Index(dataOps ...DataFunc) ([]byte, error) {

	var (
		buffer bytes.Buffer
	)

	d := &data{}

	for _, option := range dataOps {
		option(d)
	}

	if err := indexTemplate.Execute(&buffer, d); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
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
