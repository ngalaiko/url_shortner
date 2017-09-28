package template

import (
	"bytes"
	"html/template"
	"github.com/ngalayko/url_shortner/server/config"
)

const (
	dataPath = "template/data/"
	index    = "index.html"
)

var (
	indexTemplate = template.Must(template.New("index").Parse(string(MustAsset(dataPath + index))))
)

type data struct {
	FacebookApiSDK string
	FacebookAppID  string
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
		d.FacebookApiSDK = cfg.FacebookApiSDK
		d.FacebookAppID = cfg.FacebookAppID
	}
}
