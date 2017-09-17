package template

import (
	"bytes"
	"html/template"
)

const (
	dataPath = "template/data/"
	header   = "header.html"
	footer   = "footer.html"

	defaultName = "Url shorter"
)

var (
	indexTemplate = template.Must(template.New("index").
		Parse(string(MustAsset(dataPath+header)) + string(MustAsset(dataPath+footer))))
)

type data struct {
	Name string
}

// DataFunc is a func to modify template data
type DataFunc func(*data)

// Index return index.html template
func Index(dataOps ...DataFunc) ([]byte, error) {

	var (
		buffer bytes.Buffer
	)

	d := &data{
		Name: defaultName,
	}

	for _, option := range dataOps {
		option(d)
	}

	if err := indexTemplate.Execute(&buffer, d); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// WithName sets template name
func WithName(name string) DataFunc {
	return func(d *data) {
		d.Name = name
	}
}
