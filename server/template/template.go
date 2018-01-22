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
	headerFileName   = "header.html"
	footerFileName   = "footer.html"
)

var (
	headerTemplate   = template.Must(template.New("header").Parse(string(MustAsset(dataPath + headerFileName))))
	indexTemplate    = template.Must(template.New("index").Parse(string(MustAsset(dataPath + indexFileName))))
	notFoundTemplate = template.Must(template.New("notFound").Parse(string(MustAsset(dataPath + notFoundFileName))))
	footerTemplate   = template.Must(template.New("footer").Parse(string(MustAsset(dataPath + footerFileName))))
)

type data struct {
	Config config.FacebookConfig
	Errors []error

	User  *schema.User
	Links []*schema.Link
}

// DataFunc is a func to modify template data
type DataFunc func(*data)

// header returns header html
func header(d *data) ([]byte, error) {
	var headerBuffer bytes.Buffer

	if err := headerTemplate.Execute(&headerBuffer, d); err != nil {
		return nil, err
	}

	return headerBuffer.Bytes(), nil
}

// footer return footer html
func footer(d *data) ([]byte, error) {
	var footerBuffer bytes.Buffer

	if err := footerTemplate.Execute(&footerBuffer, d); err != nil {
		return nil, err
	}

	return footerBuffer.Bytes(), nil
}

// NotFound return not found page
func NotFound(dataOps ...DataFunc) ([]byte, error) {
	var notFoundBuffer bytes.Buffer

	d := parseOptions(dataOps...)

	headerBytes, err := header(d)
	if err != nil {
		return nil, err
	}

	footerBytes, err := footer(d)
	if err != nil {
		return nil, err
	}

	if err := notFoundTemplate.Execute(&notFoundBuffer, d); err != nil {
		return nil, err
	}

	return concatBytes(headerBytes, notFoundBuffer.Bytes(), footerBytes), nil
}

// Index return index page
func Index(dataOps ...DataFunc) ([]byte, error) {
	var indexBuffer bytes.Buffer

	d := parseOptions(dataOps...)

	headerBytes, err := header(d)
	if err != nil {
		return nil, err
	}

	footerBytes, err := footer(d)
	if err != nil {
		return nil, err
	}

	if err := indexTemplate.Execute(&indexBuffer, d); err != nil {
		return nil, err
	}

	return concatBytes(headerBytes, indexBuffer.Bytes(), footerBytes), nil
}

func concatBytes(bb ...[]byte) []byte {
	buffer := &bytes.Buffer{}
	for _, b := range bb {
		buffer.Write(b)
	}

	return buffer.Bytes()
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
