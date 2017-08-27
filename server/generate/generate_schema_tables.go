//go:generate go run generate_schema_tables.go
package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"text/template"

	"github.com/ngalayko/url_shortner/server/schema"
)

const (
	idTag        = "id"
	filePath     = "../dao/tables/"
	fileTemplate = `// Code generated by generate_schema_tables.go DO NOT EDIT.

package tables

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/ngalayko/url_shortner/server/dao"
	"github.com/ngalayko/url_shortner/server/helpers"
	"github.com/ngalayko/url_shortner/server/schema"
)

// Select{{ $.Name }}ById returns {{ $.Name }} from db or cache
func (t *Tables) Select{{ $.Name }}ById(id uint64) (*schema.{{ $.Name }}, error) {
	ids := []uint64{id}

	{{ alias $.Name }}{{ alias $.Name }}, err := t.Select{{ $.Name }}ByIds(ids)
	if err != nil {
		return nil, err
	}

	return {{ alias $.Name }}{{ alias $.Name }}[0], nil
}

// Select{{ $.Name }}ByIds returns {{ $.Name }}s from db or cache
func (t *Tables) Select{{ $.Name }}ByIds(ids []uint64) ([]*schema.{{ $.Name }}, error) {

	{{ alias $.Name }}{{ alias $.Name }} := make([]*schema.{{ $.Name }}, 0, len(ids))

	missingIds := make([]uint64, 0, len(ids))
	for _, id := range ids {
		value, ok := t.cache.Load(t.{{ $.TableName }}CacheKey(id))
		if !ok {
			missingIds = append(missingIds, id)
			continue
		}

		{{ alias $.Name }}{{ alias $.Name }} = append({{ alias $.Name }}{{ alias $.Name }}, value.(*schema.{{ $.Name }}))
	}

	if len(missingIds) == 0 {
		return {{ alias $.Name }}{{ alias $.Name }}, nil
	}

	{{ alias $.Name }}{{ alias $.Name }}Missing := make([]*schema.{{ $.Name }}, 0, len(missingIds))
	if err := t.db.Select(&{{ alias $.Name }}{{ alias $.Name }}Missing,
		"SELECT * "+
			"FROM {{ $.TableName }} "+
			"WHERE id IN ("+helpers.Uint64sToString(missingIds)+")",
	); err != nil {
		return nil, err
	}

	for _, {{ alias $.Name }}Missing := range {{ alias $.Name }}{{ alias $.Name }}Missing {
		{{ alias $.Name }}{{ alias $.Name }} = append({{ alias $.Name }}{{ alias $.Name }}, {{ alias $.Name }}Missing)
		t.cache.Store(t.{{ $.TableName }}CacheKey({{ alias $.Name }}Missing.ID), {{ alias $.Name }}Missing)
	}

	return {{ alias $.Name }}{{ alias $.Name }}, nil
}

// Insert{{ $.Name }} inserts {{ $.Name }} in db and cache
func (t *Tables) Insert{{ $.Name }}({{ alias $.Name }} *schema.{{ $.Name }}) error {
	return t.db.Mutate(func(tx *dao.Tx) error {

		insertSQL := "INSERT INTO {{ $.TableName }} " +
			"({{ head .DbFields}}{{ range tail .DbFields }}, {{ . }}{{ end }}) " +
			"VALUES " +
			"($1{{ range $index, $element := tail .Fields }}, ${{ sum $index 2 }}{{ end }}) " +
			"RETURNING id"

		var id uint64
		if err := tx.Get(&id, insertSQL, {{ alias $.Name }}.{{ head .Fields }}{{ range tail .Fields }}, {{ alias $.Name }}.{{ . }}{{ end }}); err != nil {
			return err
		}
		{{ alias $.Name }}.ID = id

		t.logger.Info("{{ $.Name }} created",
			zap.Reflect("$.Name", {{ alias $.Name }}),
		)
		t.cache.Store(t.{{ $.TableName }}CacheKey({{ alias $.Name }}.ID), {{ alias $.Name }})
		return nil
	})
}

// Update{{ $.Name }} updates {{ $.Name }} in db and cache
func (t *Tables) Update{{ $.Name }}({{ alias $.Name }} *schema.{{ $.Name }}) error {
	return t.db.Mutate(func(tx *dao.Tx) error {

		updateSQL := "UPDATE {{ $.TableName }} " +
			"SET " +
			{{ range $index, $element := body $.DbFields }}"{{ $element }} = ${{ sum $index 1 }}, " +
			{{ end }}"{{ last .DbFields }} = ${{ len .Fields }} " +
			fmt.Sprintf("WHERE id = %d", {{ alias $.Name }}.ID)

		_, err := tx.Exec(updateSQL, {{ alias $.Name }}.{{ head .Fields }}{{ range tail .Fields }}, {{ alias $.Name }}.{{ . }}{{ end }})
		if err != nil {
			return err
		}

		t.logger.Info("{{ $.Name }} updated",
			zap.Reflect("$.Name", {{ alias $.Name }}),
		)
		t.cache.Store(t.{{ $.TableName }}CacheKey({{ alias $.Name }}.ID), {{ alias $.Name }})
		return nil
	})
}

func (t *Tables) {{ $.TableName }}CacheKey(id uint64) string {
	return fmt.Sprintf("{{ $.Name }}:%d", id)
}
`
)

type templateData struct {
	Name         string
	TableName    string
	Fields       []string
	DbFields     []string
	UniqueFields []string
}

func main() {

	tables := []struct {
		typ       reflect.Type
		tableName string
	}{
		{reflect.TypeOf(schema.User{}), "users"},
		{reflect.TypeOf(schema.Link{}), "links"},
	}

	for _, table := range tables {
		if err := generate(table.typ, table.tableName); err != nil {
			panic(err)
		}
	}
}

func generate(typ reflect.Type, tableName string) error {

	file, err := os.Create(filePath + tableName + ".go")
	if err != nil {
		return fmt.Errorf("error opening file %s: %s", typ.Name(), err)
	}
	defer file.Close()

	data := templateData{
		Name:      typ.Name(),
		TableName: tableName,
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		dbTag := getTag(field, "db")

		if getTag(field, "unique") == "true" {
			data.UniqueFields = append(data.UniqueFields, field.Name)
		}

		if dbTag == idTag {
			continue
		}

		data.DbFields = append(data.DbFields, dbTag)
		data.Fields = append(data.Fields, field.Name)
	}

	t := template.Must(template.New(typ.Name()).Funcs(getTemplateFuncs()).Parse(fileTemplate))
	if err := t.Execute(file, data); err != nil {
		return fmt.Errorf("error executing template %s: %s", t.Name(), err)
	}

	return nil
}

func getTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"alias": func(str string) string {
			if len(str) < 1 {
				return ""
			}

			return strings.ToLower(str[:1])
		},
		"head": func(ss []string) string {
			return ss[0]
		},
		"tail": func(ss []string) []string {
			return ss[1:]
		},
		"body": func(ss []string) []string {
			if len(ss) == 1 {
				return ss
			}

			return ss[:len(ss)-1]
		},
		"last": func(ss []string) string {
			return ss[len(ss)-1]
		},
		"sum": func(a, b int) int {
			return a + b
		},
		"len": func(str []string) int {
			return len(str)
		},
	}
}

func getTag(f reflect.StructField, tagName string) string {
	tag := f.Tag.Get(tagName)

	return strings.Split(tag, ",")[0]
}
