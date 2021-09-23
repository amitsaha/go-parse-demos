package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
)

func main() {
	var files []*ast.File
	fset := token.NewFileSet()
	for _, goFile := range os.Args[1:] {
		f, err := parser.ParseFile(fset, goFile, nil, parser.ParseComments)
		if err != nil {
			log.Fatal(err)
		}
		files = append(files, f)
	}

	for _, file := range files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.FuncDecl:
				x.Body = &ast.BlockStmt{}
				fmt.Println(render(fset, x))
			case *ast.TypeSpec:
				fmt.Printf("type %s\n", render(fset, x))
			case *ast.File:
				fmt.Printf("package %s\n", x.Name)
				for _, o := range x.Imports {
					fmt.Printf("import %s\n", o.Path.Value)
				}
			}
			return true
		})
	}
}

// render returns the pretty-print of the given node
// https://arslan.io/2019/06/13/using-go-analysis-to-write-a-custom-linter/
func render(fset *token.FileSet, x interface{}) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, x); err != nil {
		panic(err)
	}
	return buf.String()
}
