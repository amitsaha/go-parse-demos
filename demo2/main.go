package main

import (
	"fmt"
	"go/ast"
	"go/parser"
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
			var fExpr ast.Expr
			switch x := n.(type) {
			case *ast.CallExpr:
				fExpr = x.Fun
			}
			if fExpr != nil {
				fmt.Printf("%#v -> %v\n", fExpr, fExpr)
			}
			return true
		})
	}
}
