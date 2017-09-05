package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"text/template"
)

var (
	structs    = []structure{}
)

const (
	schemaPath = "../schema"
)

type structure struct {
	Name   string
	Fields []string
}

type VisitorFunc func(n ast.Node) ast.Visitor

func (f VisitorFunc) Visit(n ast.Node) ast.Visitor { return f(n) }

func init() {
	if err := getSchemaNamesByAst(); err != nil {
		panic(err)
	}
}

func getSchemaNamesByAst() error {

	fs := token.NewFileSet()

	pkgs, err := parser.ParseDir(fs, schemaPath, nil, 0)
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		ast.Walk(VisitorFunc(findTypes), pkg)
	}

	return nil
}

func findTypes(n ast.Node) ast.Visitor {
	switch n := n.(type) {
	case *ast.Package:
		return VisitorFunc(findTypes)
	case *ast.File:
		return VisitorFunc(findTypes)
	case *ast.GenDecl:
		if n.Tok == token.TYPE {
			return VisitorFunc(findTypes)
		}
	case *ast.TypeSpec:
		structs = append(structs, structure{
			Name: n.Name.Name,
		})

		return VisitorFunc(findTypes)

	case *ast.StructType:

		for _, f := range n.Fields.List {
			structs[len(structs)-1].Fields = append(structs[len(structs)-1].Fields, f.Names[0].Name)
		}
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
