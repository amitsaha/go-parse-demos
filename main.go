package main

import "go/ast"
import "os"
import "go/parser"
import "log"
import "bytes"
import "go/token"
import "fmt"
import "go/printer"
import "strings"

func main() {
	var files []*ast.File

	// we will "fix" this with fmt.Println("Hi there")
	fmt.Printf("Hi there\n")
	fmt.Println("Hello there")

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
			var fExp ast.Expr
			var fArgs []ast.Expr
			switch x := n.(type) {
			case *ast.CallExpr:
				fExp = x.Fun
				fArgs = x.Args

			}
			if fExp != nil {
				for _, s := range fArgs {
					// we only care about string arguments ending with "\n"
					bl, ok := s.(*ast.BasicLit)
					if ok {										
						if strings.HasSuffix(bl.Value, `\n"`) && !strings.ContainsAny(bl.Value, "%"){
							fmt.Printf("Found: %v: %s\n", fset.Position(n.Pos()), render(fset, n))
							fmt.Printf("Recommended fix: fmt.Println(%v)\n", strings.TrimSuffix(bl.Value, `\n"`) + `"`)
						}
					}
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
