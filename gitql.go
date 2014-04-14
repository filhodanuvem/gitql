package main

import (
	"fmt"
	"flag"
	"github.com/cloudson/gitql/parser"
	"github.com/cloudson/gitql/runtime"
	"github.com/cloudson/gitql/semantical"
	"path/filepath"
)

func main() {

	query := flag.String("q", "", "The Query to search")
	pathString := flag.String("p", ".", "The (optional) path to run gitql")
	version := flag.Bool("v", false, "The version of gitql")
	flag.Parse()

	if *version {
		// @todo refactor to dynamic value
		fmt.Println("Gitql 1.0.0")
		return 
	}

	path, errFile := filepath.Abs(*pathString)

	if errFile != nil {
		panic(errFile)
	}

	parser.New(*query)
	ast, errGit := parser.AST()
	if errGit != nil {
		panic(errGit)
	}
	ast.Path = &path
	errGit = semantical.Analysis(ast)
	if errGit != nil {
		panic(errGit)
	}

	runtime.Run(ast)
}
